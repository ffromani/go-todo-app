package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/sys/unix"
)

func TestFSDirNew(t *testing.T) {
	dr := t.TempDir()
	_, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
}

func TestFSDirNoPermissions(t *testing.T) {
	if unix.Getuid() == 0 {
		t.Skipf("this test must run as non-root")
	}
	_, err := NewFSDir(filepath.Join("/usr/local", "fsdir_store"))
	if err == nil {
		t.Fatalf("directory created where it should fail")
	}
}

func TestFSDirClose(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	if err := st.Close(); err != nil {
		t.Fatalf("failed to close")
	}
}

func TestFSDirLoadAllEmpty(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	items, err := st.LoadAll()
	if err != nil {
		t.Fatalf("failed to load all: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("found items in empty store")
	}
}

func TestFSDirCreate(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	err = st.Create(ID("123"), Blob("foobar"))
	if err != nil {
		t.Fatalf("creation error: %v", err)
	}
}

func TestFSDirCreateExisting(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	err = st.Create(ID("123"), Blob("foobar"))
	if err != nil {
		t.Fatalf("creation error: %v", err)
	}
	err = st.Create(ID("123"), Blob("foobar"))
	if !errors.Is(err, ErrCorruptedContent{}) {
		t.Fatalf("create over existing entry did not fail as expected: %v", err)
	}
}

func TestFSDirSaveNonExisting(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	err = st.Save(ID("123"), Blob("foobar"))
	if !errors.Is(err, ErrCorruptedContent{}) {
		t.Fatalf("create over existing entry did not fail as expected: %v", err)
	}
}

func TestFSDirCreateSave(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	err = st.Create(ID("123"), Blob("foobar"))
	if err != nil {
		t.Fatalf("creation error: %v", err)
	}
	err = st.Save(ID("123"), Blob("fizzbuzz"))
	if err != nil {
		t.Fatalf("save existing entry failed: %v", err)
	}
}

func TestFSDirCreateLoadSaveIntermixed(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	blob, err := st.Load(ID("123"))
	if err == nil || !errors.Is(err, ErrNotFound{}) {
		t.Fatalf("load over empty unexpected error: %v blob: %v", err, blob)
	}

	err = st.Create(ID("123"), Blob("foobar"))
	if err != nil {
		t.Fatalf("creation error: %v", err)
	}
	blob, err = st.Load(ID("123"))
	if err != nil || !reflect.DeepEqual(blob, Blob("foobar")) {
		t.Fatalf("load after create unexpected error: %v blob: %v", err, blob)
	}

	err = st.Save(ID("123"), Blob("fizzbuzz"))
	if err != nil {
		t.Fatalf("save existing entry failed: %v", err)
	}
	blob, err = st.Load(ID("123"))
	if err != nil || !reflect.DeepEqual(blob, Blob("fizzbuzz")) {
		t.Fatalf("load after create unexpected error: %v blob: %v", err, blob)
	}
}

func TestFSDirCreateLoadAll(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}

	count := 9 // 9 is small random number, nothing special
	for idx := 0; idx < count; idx++ {
		id := ID(fmt.Sprintf("id#%02d", idx))
		blob := Blob(fmt.Sprintf("blob_payload_%02d", idx))
		err = st.Create(id, blob)
		if err != nil {
			t.Fatalf("creation error: %v id: %v", err, id)
		}
	}
	items, err := st.LoadAll()
	if err != nil {
		t.Fatalf("failed to load all: %v", err)
	}
	if len(items) != count {
		t.Fatalf("found mismatched items in store: got %d expected %d", len(items), count)
	}

	// random smoke test
	found := false
	for _, item := range items {
		if reflect.DeepEqual(item.ID, ID("id#05")) && reflect.DeepEqual(item.Blob, Blob("blob_payload_05")) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created item not found in loadAll result set")
	}
}

func TestFSDirDeleteEmpty(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	err = st.Delete(ID("123"))
	if !errors.Is(err, ErrNotFound{}) {
		t.Fatalf("delete unexisting entry failed unexpectedly: %v", err)
	}
}

func TestFSDirCreateDelete(t *testing.T) {
	dr := t.TempDir()
	st, err := NewFSDir(filepath.Join(dr, "fsdir_store"))
	if err != nil {
		t.Fatalf("failed to create")
	}
	blob, err := st.Load(ID("123"))
	if err == nil || !errors.Is(err, ErrNotFound{}) {
		t.Fatalf("load over empty unexpected error: %v blob: %v", err, blob)
	}

	err = st.Create(ID("123"), Blob("foobar"))
	if err != nil {
		t.Fatalf("creation error: %v", err)
	}
	blob, err = st.Load(ID("123"))
	if err != nil || !reflect.DeepEqual(blob, Blob("foobar")) {
		t.Fatalf("load after create unexpected error: %v blob: %v", err, blob)
	}

	err = st.Delete(ID("123"))
	if err != nil {
		t.Fatalf("delete existing entry failed: %v", err)
	}
	blob, err = st.Load(ID("123"))
	if err == nil || !errors.Is(err, ErrNotFound{}) {
		t.Fatalf("load over empty unexpected error: %v blob: %v", err, blob)
	}
}
