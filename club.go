package computerclub

import "time"

type Club struct {
	NumberOfTables int
	StartTime      time.Time
	EndTime        time.Time
	Coast          int
}

type ActiveClientCard struct {
	StartTime   time.Time
	MiddleTime  time.Time
	EntTime     time.Time
	TableNumber int
}

type DayBalance struct {
	Money     int
	Time      time.Time
	MoneyTime int
}
