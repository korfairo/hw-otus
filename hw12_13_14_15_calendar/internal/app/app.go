package app

import (
	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/logger"
	"time"

	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage Storage
	logger  Logger
}

type Logger interface {
	WithField(key string, value interface{}) *logger.Logger
	WithError(err error) *logger.Logger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}

type Storage interface {
	AddEvent(event storage.Event) error
	UpdateEvent(id uint32, event storage.Event) error
	DeleteEvent(id uint32) error
	GetEvent(id uint32) (storage.Event, error)
	GetEventsOnDay(t time.Time) ([]storage.Event, error)
	GetEventsOnWeek(t time.Time) ([]storage.Event, error)
	GetEventsOnMonth(t time.Time) ([]storage.Event, error)
}

func New(storage Storage, logger Logger) *App {
	return &App{
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(id uint32, title string) error {
	return a.storage.AddEvent(storage.Event{
		ID:    id,
		Title: title,
	})
}
