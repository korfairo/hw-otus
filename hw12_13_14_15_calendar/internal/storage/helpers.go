package storage

import "time"

func GetDayInterval(t time.Time) (from, to time.Time) {
	year, month, day := t.Date()
	from = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	to = from.Add(24 * time.Hour)
	return
}

func GetWeekIntervals(t time.Time) (from, to time.Time) {
	year, month, day := t.Date()
	dayStart := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	weekday := int(dayStart.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	from = dayStart.AddDate(0, 0, -1*weekday+1)
	to = from.Add(7 * 24 * time.Hour)
	return
}

func GetMonthInterval(t time.Time) (from, to time.Time) {
	year, month, _ := t.Date()
	from = time.Date(year, month, 0, 0, 0, 0, 0, t.Location())
	to = from.AddDate(0, 1, 0)
	return
}
