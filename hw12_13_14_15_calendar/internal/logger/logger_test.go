package logger

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerLevels(t *testing.T) {
	logAllLevels := func(l *Logger) {
		l.Debug("debug")
		l.Debugf("%s", "debug")

		l.Info("info")
		l.Infof("%s", "info")

		l.Warn("warn")
		l.Warnf("%s", "warn")

		l.Error("error")
		l.Errorf("%s", "error")
	}

	type args struct {
		level string
	}

	type expected struct {
		wantErr bool
		rowsNum int
	}
	tests := []struct {
		name     string
		args     args
		action   func(l *Logger)
		expected expected
	}{
		{
			name: "unknown log level",
			args: args{
				level: "random",
			},
			expected: expected{
				wantErr: true,
			},
		},
		{
			name: "debug level",
			args: args{
				level: "debug",
			},
			action: logAllLevels,
			expected: expected{
				rowsNum: 8,
			},
		},
		{
			name: "info level",
			args: args{
				level: "info",
			},
			action: logAllLevels,
			expected: expected{
				rowsNum: 6,
			},
		},
		{
			name: "warn level",
			args: args{
				level: "warn",
			},
			action: logAllLevels,
			expected: expected{
				rowsNum: 4,
			},
		},
		{
			name: "error level",
			args: args{
				level: "error",
			},
			action: logAllLevels,
			expected: expected{
				rowsNum: 2,
			},
		},
		{
			name: "fatal level",
			args: args{
				level: "fatal",
			},
			action: logAllLevels,
			expected: expected{
				rowsNum: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPaths := createTempFiles(t, "1")

			logger, err := New(Config{
				Level:       tt.args.level,
				OutputPaths: outputPaths,
			})
			if tt.expected.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.action != nil {
				tt.action(logger)

				outputData := readFile(t, outputPaths[0])
				gotRowsNum := bytes.Count(outputData, []byte("\n"))

				assert.Equal(t, tt.expected.rowsNum, gotRowsNum, "the number of rows written is different from the expected number")
			}
		})
	}
}

func TestPanic(t *testing.T) {
	logger, err := New(Config{
		Level: "info",
	})
	assert.NoError(t, err)

	assert.Panics(t, func() {
		logger.Panic("test panic #1")
	}, "logger.Panic() must panic")

	assert.Panics(t, func() {
		logger.Panicf("test panic #%d", 2)
	}, "logger.Panicf() must panic")
}

func TestWriteToFile(t *testing.T) {
	files := createTempFiles(t, "1.log", "2.log")

	t.Logf("use temp files: %v", strings.Join(files, ", "))

	log, err := New(Config{
		Level:       "info",
		OutputPaths: files,
	})
	assert.NoError(t, err)
	defer func() {
		err = log.Sync()
		assert.NoError(t, err)
	}()

	log.Info("test data", 1)

	for _, filename := range files {
		data := readFile(t, filename)
		if len(data) == 0 {
			t.Errorf("no data in file %s", filename)
		}
	}
}

func createTempFiles(t *testing.T, filenames ...string) (filepaths []string) {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})

	t.Logf("use temp dir: %s", dir)

	for _, name := range filenames {
		filepaths = append(filepaths, path.Join(dir, name))
	}
	return filepaths
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()

	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
