package storage

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uint32 `db:"uuid"`
	UserID      uint32 `db:"user_uuid"`
	Title       string `db:"title"`
	Description string `db:"description"`

	Date     time.Time     `db:"date"`
	Duration time.Duration `db:"duration"`
	Notice   time.Duration `db:"notice"`
}

func NewUUID() uint32 {
	return uuid.New().ID()
}
