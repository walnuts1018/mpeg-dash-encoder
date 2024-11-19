package temp

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	type args struct {
		r        io.Reader
		filename string
	}
	tests := []struct {
		name        string
		args        args
		wantPath    string
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "normal",
			args: args{
				r:        strings.NewReader("test"),
				filename: "TestCreateTempFile",
			},
			wantContent: []byte("test"),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTempFile(tt.args.r, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()

			gotContent, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("failed to read all: %v", err)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("CreateTempFile() = %v, want %v", string(gotContent), string(tt.wantContent))
				return
			}
		})
	}
}
