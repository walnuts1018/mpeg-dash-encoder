package temp

import (
	"io"
	"os"
)

func CreateTempFile(r io.Reader, filename string) (*os.File, error) {
	inputFile, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(inputFile, r)
	if err != nil {
		return nil, err
	}
	if err := inputFile.Sync(); err != nil {
		return nil, err
	}
	if _, err := inputFile.Seek(0, 0); err != nil {
		return nil, err
	}

	return inputFile, nil
}
