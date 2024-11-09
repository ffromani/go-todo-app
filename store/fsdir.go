package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FSDir use a filesystem directory as backing store, and one-file-per-blob
// approach in that directory.
// Known limitations:
// - no directory escape checks
// - no protection against concurrent access between goroutines
// - no protection against concurrent access between different processes
type FSDir struct {
	fsPath string
}

var _ Storage = &FSDir{}

func NewFSDir(fsPath string) (*FSDir, error) {
	fsDir := FSDir{
		fsPath: fsPath,
	}
	err := fsDir.ensure()
	return &fsDir, err
}

func (dr *FSDir) Close() error {
	return nil // nothing to do
}

func (dr *FSDir) LoadAll() ([]Item, error) {
	var items []Item
	err := filepath.WalkDir(dr.fsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == dr.fsPath {
			return nil // skip self when iterating, see filepath.WalkDir doscs
		}
		fName := filepath.Base(path)
		if strings.HasPrefix(fName, ".") {
			return nil // ignore
		}
		if d.IsDir() {
			return ErrCorruptedContent{
				Name:     path,
				IntError: fmt.Errorf("%q is a directory", path),
			}
		}
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			return nil // ignore
		}
		if err != nil {
			return ErrCorruptedContent{
				Name:     path,
				IntError: err,
			}
		}
		items = append(items, Item{
			ID:   ID(fName),
			Blob: Blob(data),
		})
		return nil
	})
	return items, err
}

func (dr *FSDir) Load(id ID) (Blob, error) {
	objPath := filepath.Join(dr.fsPath, string(id))
	data, err := os.ReadFile(objPath)
	if os.IsNotExist(err) {
		return nil, ErrNotFound{ID: id}
	}
	return Blob(data), err
}

func (dr *FSDir) Create(id ID, blob Blob) error {
	objPath := filepath.Join(dr.fsPath, string(id))
	if _, err := os.Stat(objPath); err == nil {
		return ErrCorruptedContent{Name: objPath}
	}
	return os.WriteFile(objPath, blob, 0644)
}

func (dr *FSDir) Save(id ID, blob Blob) error {
	objPath := filepath.Join(dr.fsPath, string(id))
	if _, err := os.Stat(objPath); err != nil {
		return ErrCorruptedContent{Name: objPath}
	}
	return os.WriteFile(objPath, blob, 0644)
}

func (dr *FSDir) Delete(id ID) error {
	objPath := filepath.Join(dr.fsPath, string(id))
	err := os.Remove(objPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound{ID: id}
		}
		return ErrCorruptedContent{
			Name:     objPath,
			IntError: err,
		}
	}
	return nil
}

func (dr *FSDir) ensure() error {
	return os.MkdirAll(dr.fsPath, 0750)
}
