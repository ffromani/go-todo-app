package store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	lockFile string = ".lock"
)

type ErrDifferentOwner struct {
	pid int
}

func (e ErrDifferentOwner) Error() string {
	return fmt.Sprintf("owned by pid %d", e.pid)
}

type FSDir struct {
	lastObjectID ID
	pid          int
	fsPath       string
}

func NewFSDir(fsPath string) (*FSDir, error) {
	fsDir := FSDir{
		pid:    os.Getpid(),
		fsPath: fsPath,
	}
	if err := fsDir.checkOwnedByMe(); err != nil {
		return nil, err
	}
	lastObjectID, err := fsDir.getLastObjectID()
	if err != nil {
		return nil, err
	}
	fsDir.lastObjectID = max(lastObjectID, 1)
	return &fsDir, nil
}

func (dr *FSDir) getLastObjectID() (ID, error) {
	var lastObjectID ID
	rerr := filepath.WalkDir(dr.fsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return ErrCorruptedContent{Name: path}
		}
		fName := filepath.Base(path)
		if fName == lockFile {
			return nil // treat explicitely our metadata
		}
		if strings.HasPrefix(fName, ".") {
			return nil // ignore
		}
		objectID, cerr := strconv.Atoi(fName)
		if cerr != nil {
			return ErrCorruptedContent{Name: path}
		}
		lastObjectID = max(ID(objectID), lastObjectID)
		return nil
	})
	return lastObjectID, rerr
}

func (dr *FSDir) getOwner() (int, error) {
	lockPath := filepath.Join(dr.fsPath, lockFile)
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func (dr *FSDir) checkOwnedByMe() error {
	curPid, err := dr.getOwner()
	if err != nil {
		return err
	}
	if curPid != dr.pid {
		return ErrDifferentOwner{pid: curPid}
	}
	return nil
}

func (dr *FSDir) setMeAsOwner() error {
	lockPath := filepath.Join(dr.fsPath, lockFile)
	return os.WriteFile(lockPath, []byte(strconv.Itoa(dr.pid)), 0644)
}

func (dr *FSDir) Close() error {
	return nil
}

func (dr *FSDir) Create(data Blob) (ID, error) {
	if err := dr.checkOwnedByMe(); err != nil {
		return 0, err
	}
	objectID := dr.lastObjectID + 1
	err := dr.Save(objectID, data)
	if err != nil {
		return NullID, err
	}
	dr.lastObjectID = objectID
	return dr.lastObjectID, nil
}

func (dr *FSDir) LoadAll() ([]Item, error) {
	if err := dr.checkOwnedByMe(); err != nil {
		return nil, err
	}

	var items []Item
	err := filepath.WalkDir(dr.fsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return ErrCorruptedContent{Name: path}
		}
		fName := filepath.Base(path)
		if fName == lockFile {
			return nil // treat explicitely our metadata
		}
		if strings.HasPrefix(fName, ".") {
			return nil // ignore
		}
		objectID, cerr := strconv.Atoi(fName)
		if cerr != nil {
			return ErrCorruptedContent{Name: path}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return ErrCorruptedContent{Name: path}
		}
		items = append(items, Item{
			ID:   ID(objectID),
			Blob: Blob(data),
		})
		return nil
	})
	return items, err
}

func (dr *FSDir) Load(id ID) (Blob, error) {
	if err := dr.checkOwnedByMe(); err != nil {
		return nil, err
	}
	objPath := filepath.Join(dr.fsPath, strconv.FormatInt(int64(id), 10))
	data, err := os.ReadFile(objPath)
	if os.IsNotExist(err) {
		return nil, ErrNotFound{ID: id}
	}
	return Blob(data), err
}

func (dr *FSDir) Save(id ID, blob Blob) error {
	if err := dr.checkOwnedByMe(); err != nil {
		return err
	}
	objPath := filepath.Join(dr.fsPath, strconv.FormatInt(int64(id), 10))
	return os.WriteFile(objPath, blob, 0644)
}

func (dr *FSDir) Delete(id ID) error {
	if err := dr.checkOwnedByMe(); err != nil {
		return err
	}
	objPath := filepath.Join(dr.fsPath, strconv.FormatInt(int64(id), 10))
	return os.Remove(objPath)
}
