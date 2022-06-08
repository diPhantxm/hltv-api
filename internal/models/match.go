package models

import "time"

type Match struct {
	Id               int       `json:"Id"`
	TeamA            Team      `json:"teamA"`
	TeamB            Team      `json:"teamB"`
	StartTime        time.Time `json:"start"`
	Maps             []string  `json:"maps,omitempty"`
	Viewers          int       `json:"viewers"`
	PlayerOfTheMatch Player    `json:"playerOfTheMatch,omitempty"`
}
