package models

type Team struct {
	Id           int           `json:"id"`
	Ranking      int           `json:"ranking"`
	WeeksInTop30 int           `json:"weeksInTop30"`
	Name         string        `json:"name"`
	Country      string        `json:"country"`
	Social       []Social      `json:"social,omitempty"`
	Achievements []Achievement `json:"achievements,omitempty"`
	AverageAge   float32       `json:"averageAge"`
}
