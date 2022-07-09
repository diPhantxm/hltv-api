package models

import "time"

type Match struct {
	Id               int       `json:"Id"`
	TeamA            string    `json:"teamA"`
	TeamB            string    `json:"teamB"`
	StartTime        time.Time `json:"start"`
	Maps             []Map     `json:"maps"`
	Viewers          int       `json:"viewers"`
	PlayerOfTheMatch string    `json:"playerOfTheMatch,omitempty"`
	IsOver           bool      `json:"isOver"`
}
