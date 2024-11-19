package hash

import (
	"io"
	"strings"
	"testing"
)

func TestReaderHash(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				r: strings.NewReader("testtest"),
			},
			want:    "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578", // echo -n 'testtest' | shasum -a 256
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				r: strings.NewReader(""),
			},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // echo -n '' | shasum -a 256
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReaderHash(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReaderHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReaderHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
