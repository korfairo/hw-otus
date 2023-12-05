package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.SugaredLogger
}

type Config struct {
	Level       string
	OutputPaths []string
}

func New(cfg Config) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal log level")
	}

	if cfg.OutputPaths == nil {
		cfg.OutputPaths = []string{"stdout"}
	}

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:       cfg.OutputPaths,
		Development:       true,
		DisableStacktrace: true,
		Encoding:          "console",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
	}

	z, err := zapCfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build logger")
	}

	return &Logger{
		zap: z.Sugar(),
	}, nil
}

// WithField returns a cloned logger with a new field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.withFields(zap.Any(key, value))
}

// WithError is a shorthand for Logger.WithField("error", err)
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err)
}

func (l *Logger) Debug(args ...interface{})                 { l.zap.Debug(args...) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.zap.Debugf(format, args...) }

func (l *Logger) Info(args ...interface{})                 { l.zap.Info(args...) }
func (l *Logger) Infof(format string, args ...interface{}) { l.zap.Infof(format, args...) }

func (l *Logger) Warn(args ...interface{})                 { l.zap.Warn(args...) }
func (l *Logger) Warnf(format string, args ...interface{}) { l.zap.Warnf(format, args...) }

func (l *Logger) Error(args ...interface{})                 { l.zap.Error(args...) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.zap.Errorf(format, args...) }

func (l *Logger) Fatal(args ...interface{})                 { l.zap.Fatal(args...) }
func (l *Logger) Fatalf(format string, args ...interface{}) { l.zap.Fatalf(format, args...) }

func (l *Logger) Panic(args ...interface{})                 { l.zap.Panic(args...) }
func (l *Logger) Panicf(format string, args ...interface{}) { l.zap.Panicf(format, args...) }

// Sync clears all buffered log entries before exiting the application.
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func (l *Logger) withFields(fields ...zap.Field) *Logger {
	return &Logger{
		zap: l.zap.Desugar().With(fields...).Sugar(),
	}
}
