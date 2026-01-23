package db

import (
	"fmt"
	"os"
	"path/filepath"
)

type byteStorage struct {
	storeDir string
}

func newByteStorage(storeDir string) (*byteStorage, error) {
	return &byteStorage{storeDir: storeDir}, os.MkdirAll(storeDir, 0744)
}

func (s byteStorage) filePath(id, fileType string) string {
	return filepath.Join(s.storeDir, fmt.Sprintf("%s.%s", id, fileType))
}

func (s byteStorage) Add(id, fileType string, content []byte) error {
	return os.WriteFile(s.filePath(id, fileType), content, 0744)
}

func (s byteStorage) Del(id, fileType string) error {
	return os.Remove(s.filePath(id, fileType))
}

func (s byteStorage) Load(id, fileType string) ([]byte, error) {
	return os.ReadFile(s.filePath(id, fileType))
}
