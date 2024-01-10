package memorystorage

import (
	"sort"
	"testing"
	"time"

	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/storage"

	"github.com/stretchr/testify/assert"
)

var (
	anniversaryEvent = storage.Event{
		ID:          storage.NewUUID(),
		UserID:      storage.NewUUID(),
		Title:       "Anniversary",
		Description: "Our wedding anniversary",
		Date:        mustParseRFC3339("2024-01-19T10:00:00+05:00"),
		Duration:    12 * time.Hour,
		Notice:      24 * time.Hour,
	}
	birthdayEvent = storage.Event{
		ID:          storage.NewUUID(),
		UserID:      storage.NewUUID(),
		Title:       "Birthday",
		Description: "My happy birthday",
		Date:        mustParseRFC3339("2024-01-19T12:00:00+05:00"),
		Duration:    24 * time.Hour,
		Notice:      12 * time.Hour,
	}
	newYearEvent = storage.Event{
		ID:          storage.NewUUID(),
		UserID:      storage.NewUUID(),
		Title:       "New Year",
		Description: "Winter family holiday",
		Date:        mustParseRFC3339("2024-01-01T00:00:00+03:00"),
		Duration:    2 * time.Hour,
		Notice:      24 * time.Hour,
	}
	christmasEvent = storage.Event{
		ID:          storage.NewUUID(),
		UserID:      storage.NewUUID(),
		Title:       "Christmas",
		Description: "Merry Christmas!",
		Date:        mustParseRFC3339("2024-01-07T00:00:00+03:00"),
		Duration:    24 * time.Hour,
		Notice:      48 * time.Hour,
	}
	womenDayEvent = storage.Event{
		ID:          storage.NewUUID(),
		UserID:      storage.NewUUID(),
		Title:       "8 March",
		Description: "International Women's Day",
		Date:        mustParseRFC3339("2024-03-08T00:00:00+03:00"),
		Duration:    24 * time.Hour,
		Notice:      24 * time.Hour,
	}
)

func TestStorage(t *testing.T) {
	repo := New()
	err := repo.AddEvent(birthdayEvent)
	assert.NoError(t, err)
	err = repo.AddEvent(newYearEvent)
	assert.NoError(t, err)
	err = repo.AddEvent(christmasEvent)
	assert.NoError(t, err)
	err = repo.AddEvent(womenDayEvent)
	assert.NoError(t, err)

	gotBirthday, err := repo.GetEvent(birthdayEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, birthdayEvent, gotBirthday)

	gotNewYear, err := repo.GetEvent(newYearEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, newYearEvent, gotNewYear)

	gotChristmas, err := repo.GetEvent(christmasEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, christmasEvent, gotChristmas)

	gotWomensDay, err := repo.GetEvent(womenDayEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, womenDayEvent, gotWomensDay)
}

func TestStorage_UpdateEvent(t *testing.T) {
	repo := New()
	err := repo.AddEvent(birthdayEvent)
	assert.NoError(t, err)

	newBirthdayEvent := birthdayEvent
	newTitle := "It's my Birthday"
	newNotice := 48 * time.Hour
	newBirthdayEvent.Notice = newNotice
	newBirthdayEvent.Title = newTitle

	err = repo.UpdateEvent(birthdayEvent.ID, newBirthdayEvent)
	assert.NoError(t, err)

	gotEventAfterUpdate, err := repo.GetEvent(birthdayEvent.ID)
	assert.NoError(t, err)
	assert.Equal(t, newBirthdayEvent, gotEventAfterUpdate)
}

func TestStorage_DeleteEvent(t *testing.T) {
	repo := New()
	err := repo.AddEvent(birthdayEvent)
	assert.NoError(t, err)

	err = repo.DeleteEvent(birthdayEvent.ID)
	assert.NoError(t, err)

	gotEventAfterDelete, err := repo.GetEvent(birthdayEvent.ID)
	assert.ErrorIs(t, err, ErrEventNotFound)
	assert.Equal(t, storage.Event{}, gotEventAfterDelete)
}

func TestStorage_GetEventsOnDay(t *testing.T) {
	type args struct {
		t time.Time
	}

	tests := []struct {
		name string
		args args
		want []storage.Event
	}{
		{
			name: "one event, positive test",
			args: args{
				t: mustParseRFC3339("2024-01-01T12:00:00+03:00"),
			},
			want: []storage.Event{newYearEvent},
		},
		{
			name: "two events, positive test",
			args: args{
				t: mustParseRFC3339("2024-01-19T12:00:00+05:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent},
		},
		{
			name: "positive, leftmost local time",
			args: args{
				t: mustParseRFC3339("2024-01-19T00:00:00+05:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent},
		},
		{
			name: "positive, rightmost local time",
			args: args{
				t: mustParseRFC3339("2024-01-19T23:59:59+05:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent},
		},
		{
			name: "positive, another time zone, leftmost time",
			args: args{
				t: mustParseRFC3339("2024-01-19T00:00:00+12:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent},
		},
		{
			name: "positive, another time zone, rightmost time",
			args: args{
				t: mustParseRFC3339("2024-01-19T23:59:59+12:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent},
		},
		{
			name: "negative, previous day, local time",
			args: args{
				t: mustParseRFC3339("2024-01-18T23:59:59+05:00"),
			},
			want: nil,
		},
		{
			name: "negative, next day, local time",
			args: args{
				t: mustParseRFC3339("2024-01-20T00:00:00+05:00"),
			},
			want: nil,
		},
		{
			name: "negative, previous day, another time zone",
			args: args{
				t: mustParseRFC3339("2024-01-18T23:59:59+00:00"),
			},
			want: nil,
		},
		{
			name: "negative, next day, another time zone",
			args: args{
				t: mustParseRFC3339("2024-01-20T00:00:00-12:00"),
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := New()

			err := repo.AddEvent(birthdayEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(anniversaryEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(newYearEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(christmasEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(womenDayEvent)
			assert.NoError(t, err)

			gotEvents, err := repo.GetEventsOnDay(test.args.t)
			assert.NoError(t, err)
			sort.Slice(gotEvents, func(i, j int) bool {
				return gotEvents[i].Title < gotEvents[j].Title
			})
			assert.Equal(t, test.want, gotEvents)
		})
	}
}

func TestStorage_GetEventsOnWeek(t *testing.T) {
	type args struct {
		t time.Time
	}

	tests := []struct {
		name string
		args args
		want []storage.Event
	}{
		{
			name: "one event, positive test",
			args: args{
				t: mustParseRFC3339("2024-03-05T12:00:00+03:00"),
			},
			want: []storage.Event{womenDayEvent},
		},
		{
			name: "one event, negative test",
			args: args{
				t: mustParseRFC3339("2024-03-03T12:00:00+03:00"),
			},
			want: nil,
		},
		{
			name: "two events, positive test",
			args: args{
				t: mustParseRFC3339("2024-01-01T12:00:00+03:00"),
			},
			want: []storage.Event{christmasEvent, newYearEvent},
		},
		{
			name: "two events, negative test",
			args: args{
				t: mustParseRFC3339("2023-12-31T12:00:00+03:00"),
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := New()

			err := repo.AddEvent(birthdayEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(anniversaryEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(newYearEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(christmasEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(womenDayEvent)
			assert.NoError(t, err)

			gotEvents, err := repo.GetEventsOnWeek(test.args.t)
			assert.NoError(t, err)
			sort.Slice(gotEvents, func(i, j int) bool {
				return gotEvents[i].Title < gotEvents[j].Title
			})
			assert.Equal(t, test.want, gotEvents)
		})
	}
}

func TestStorage_GetEventsOnMonth(t *testing.T) {
	type args struct {
		t time.Time
	}

	tests := []struct {
		name string
		args args
		want []storage.Event
	}{
		{
			name: "no events",
			args: args{
				t: mustParseRFC3339("2024-04-15T12:00:00+03:00"),
			},
			want: nil,
		},
		{
			name: "one event",
			args: args{
				t: mustParseRFC3339("2024-03-13T12:00:00+03:00"),
			},
			want: []storage.Event{womenDayEvent},
		},
		{
			name: "few events",
			args: args{
				t: mustParseRFC3339("2024-01-27T12:00:00+03:00"),
			},
			want: []storage.Event{anniversaryEvent, birthdayEvent, christmasEvent, newYearEvent},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := New()

			err := repo.AddEvent(birthdayEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(anniversaryEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(newYearEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(christmasEvent)
			assert.NoError(t, err)
			err = repo.AddEvent(womenDayEvent)
			assert.NoError(t, err)

			gotEvents, err := repo.GetEventsOnMonth(test.args.t)
			assert.NoError(t, err)
			sort.Slice(gotEvents, func(i, j int) bool {
				return gotEvents[i].Title < gotEvents[j].Title
			})
			assert.Equal(t, test.want, gotEvents)
		})
	}
}

// mustParseRFC3339 parses time in RFC3339 format: '2006-01-02T15:04:05Z07:00'.
// It panics if time.Parse returns an error
func mustParseRFC3339(t string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return parsedTime
}
