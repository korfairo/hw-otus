package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
}

func New(DSN string) *Storage {
	return &Storage{
		dsn: DSN,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", s.dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(event storage.Event) error {
	q := `INSERT INTO events (uuid, user_uuid, title, description, date, duration, notice) 
		  VALUES (:uuid, :user_uuid, :title, :description, :date, :duration, :notice)`
	if _, err := s.db.NamedExec(q, &event); err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateEvent(id uint32, event storage.Event) error {
	q := fmt.Sprintf(`UPDATE events
		  SET uuid=:uuid, user_uuid=:user_uuid, title=:title, description=:description, date=:date, duration=:duration, notice=:notice
		  WHERE uuid = %d`, id)
	if _, err := s.db.NamedExec(q, &event); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(id uint32) error {
	q := "DELETE FROM events WHERE uuid = $1"
	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEvent(id uint32) (e storage.Event, err error) {
	q := `SELECT uuid, user_uuid, title, description, date, duration, notice FROM events WHERE uuid = $1`
	err = s.db.Select(&e, q, id)
	return
}

type period uint8

const (
	periodDay period = iota
	periodWeek
	periodMonth
)

func (s *Storage) GetEventsOnDay(t time.Time) ([]storage.Event, error) {
	return s.getEvents(periodDay, t)
}

func (s *Storage) GetEventsOnWeek(t time.Time) ([]storage.Event, error) {
	return s.getEvents(periodWeek, t)
}

func (s *Storage) GetEventsOnMonth(t time.Time) ([]storage.Event, error) {
	return s.getEvents(periodMonth, t)
}

func (s *Storage) getEvents(p period, t time.Time) ([]storage.Event, error) {
	var events []storage.Event
	var from, to time.Time
	switch p {
	case periodDay:
		from, to = storage.GetDayInterval(t)
	case periodWeek:
		from, to = storage.GetWeekIntervals(t)
	case periodMonth:
		from, to = storage.GetMonthInterval(t)
	}

	q := `SELECT uuid, user_uuid, title, description, date, duration, notice
		  FROM events
		  WHERE date BETWEEN $1 AND $2`
	if err := s.db.Select(&events, q, from, to); err != nil {
		return nil, err
	}
	return events, nil
}
