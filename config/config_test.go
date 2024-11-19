package config

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	"dario.cat/mergo"
	_ "github.com/joho/godotenv/autoload"
)

var requiredEnvs = map[string]string{
	"MINIO_ACCESS_KEY": "test",
	"MINIO_SECRET_KEY": "test",
	"REDIS_PASSWORD":   "test",
	"ADMIN_TOKEN":      "test",
	"JWT_SIGN_SECRET":  "test",
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string //env
		want    Config
		wantErr bool
	}{
		{
			name: "check custom type default",
			envs: map[string]string{},
			want: Config{
				ServerPort: "8080",
				LogLevel:   slog.LevelInfo,
			},
			wantErr: false,
		},
		{
			name: "normal",
			envs: map[string]string{
				"SERVER_PORT": "9000",
			},
			want: Config{
				ServerPort: "9000",
			},
			wantErr: false,
		},
		{
			name: "check custom type",
			envs: map[string]string{
				"LOG_LEVEL": "debug",
			},
			want: Config{
				LogLevel: slog.LevelDebug,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			var envs = requiredEnvs
			for k, v := range tt.envs {
				envs[k] = v
			}

			for k, v := range envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set env: %v", err)
					return
				}
				defer os.Unsetenv(k)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ok, err := equal(got, tt.want)
			if err != nil {
				t.Errorf("failed to check config: %v", err)
				return
			}
			if !ok {
				t.Errorf("Load() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func equal(got, want Config) (bool, error) {
	merged := want
	if err := mergo.Merge(&merged, got); err != nil {
		return false, err
	}

	return reflect.DeepEqual(merged, got), nil
}

func Test_equal(t *testing.T) {
	type args struct {
		got  Config
		want Config
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				got: Config{
					ServerPort: "8080",
					LogLevel:   slog.LevelDebug,
				},
				want: Config{
					ServerPort: "8080",
					LogLevel:   slog.LevelDebug,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not equal",
			args: args{
				got: Config{
					ServerPort: "8080",
					LogLevel:   slog.LevelInfo,
				},
				want: Config{
					ServerPort: "9090",
					LogLevel:   slog.LevelDebug,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := equal(tt.args.got, tt.args.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("equal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
