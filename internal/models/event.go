package models

import "time"

type Event struct {
	StartDate time.Time
	EndDate   time.Time
	PrizePool int
	Teams     []Team
	Location  string
}
