package fileutil

import (
	"os"
	"path/filepath"
)

func CreateFileRecursive(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(path)
}
