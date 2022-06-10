package models

import "time"

type Match struct {
	Id               int       `json:"Id"`
	TeamA            string    `json:"teamA"`
	TeamB            string    `json:"teamB"`
	StartTime        time.Time `json:"start"`
	Maps             []string  `json:"maps,omitempty"`
	Viewers          int       `json:"viewers"`
	PlayerOfTheMatch string    `json:"playerOfTheMatch,omitempty"`
}
