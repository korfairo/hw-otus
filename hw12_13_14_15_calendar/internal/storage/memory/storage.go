package memorystorage

import (
	"sync"
	"time"

	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventNotFound      = errors.New("event not found")
)

type Storage struct {
	mu     sync.RWMutex
	events map[uint32]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[uint32]storage.Event),
	}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.eventsExists(event.ID) {
		return ErrEventAlreadyExists
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(id uint32, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.eventsExists(event.ID) {
		return ErrEventNotFound
	}
	s.events[id] = event
	return nil
}

func (s *Storage) DeleteEvent(id uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.eventsExists(id) {
		return ErrEventNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvent(id uint32) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.events[id]
	if !exists {
		return storage.Event{}, ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) eventsExists(id uint32) bool {
	_, exists := s.events[id]
	return exists
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

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.events {
		if (event.Date.After(from) || event.Date == from) && event.Date.Before(to) {
			events = append(events, event)
		}
	}
	return events, nil
}
