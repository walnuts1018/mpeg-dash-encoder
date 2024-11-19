package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileRecursive(t *testing.T) {
	baseDir := filepath.Join(os.TempDir(), "fileutil_test")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "depth 1",
			args: args{
				path: filepath.Join(baseDir, "test.txt"),
			},
			wantErr: false,
		},
		{
			name: "depth 2",
			args: args{
				path: filepath.Join(baseDir, "dir1", "test.txt"),
			},
			wantErr: false,
		},
		{
			name: "depth 4",
			args: args{
				path: filepath.Join(baseDir, "dir1", "dir2", "dir3", "test.txt"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("os.RemoveAll() error = %v", err)
			}

			got, err := CreateFileRecursive(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()

			if _, err := os.Stat(tt.args.path); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					t.Errorf("CreateFileRecursive() file not created: %v", err)
				}
				t.Errorf("os.Stat() error = %v", err)
			}
		})
	}
}
