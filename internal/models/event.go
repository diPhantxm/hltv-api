package models

import "time"

type Event struct {
	Id        int       `json:"Id"`
	Name      string    `json:"Name"`
	StartDate time.Time `json:"Start Date"`
	EndDate   time.Time `json:"End Date"`
	PrizePool int       `json:"Prize Pool"`
	Teams     []string  `json:"Teams"`
	Location  string    `json:"Location"`
}
