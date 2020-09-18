package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

var (
	ChairDraftFiles  *FileIterator
	EstateDraftFiles *FileIterator
)

type FileIterator struct {
	parentDirPath string
	files         []os.FileInfo
	offset        int
	mu            sync.Mutex
}

func NewFileIterator(parentDirPath string) (*FileIterator, error) {
	files, err := ioutil.ReadDir(parentDirPath)
	if err != nil {
		return nil, err
	}
	return &FileIterator{
		parentDirPath: parentDirPath,
		files:         files,
		offset:        0,
	}, nil
}

func (d *FileIterator) Next() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.offset >= len(d.files) {
		return "", fmt.Errorf("too many read")
	}
	f := d.files[d.offset]
	d.offset++
	return path.Join(d.parentDirPath, f.Name()), nil
}

func loadChairDraftFiles(dataDir string) error {
	files, err := NewFileIterator(path.Join(dataDir, "result/draft_data/chair"))
	if err != nil {
		return err
	}

	ChairDraftFiles = files
	return nil
}

func loadEstateDraftFiles(dataDir string) error {
	files, err := NewFileIterator(path.Join(dataDir, "result/draft_data/estate"))
	if err != nil {
		return err
	}

	EstateDraftFiles = files
	return nil
}
