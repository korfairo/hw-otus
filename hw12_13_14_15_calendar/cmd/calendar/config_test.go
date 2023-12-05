package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		configFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				"../../configs/config.toml",
			},
			want: Config{
				Logger: LoggerConf{
					Level: "INFO",
				},
				Storage: StorageConf{
					Memory:      true,
					PostgresDSN: "postgresql://user:secret@localhost:5432/calendar?sslmode=disable",
				},
				Server: ServerConf{
					Host: "localhost",
					Port: "5000",
				},
			},
			wantErr: false,
		},
		{
			name: "incorrect config path",
			args: args{
				"../../configs/unknown_config.toml",
			},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.configFilePath)
			assert.Equal(t, tt.want, got)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
