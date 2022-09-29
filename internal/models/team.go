package models

type Team struct {
	Id           int           `json:"id"`
	Ranking      int           `json:"ranking"`
	WeeksInTop30 int           `json:"weeksInTop30"`
	Name         string        `json:"name"`
	Country      string        `json:"country"`
	Roaster      Roaster       `json:"roaster"`
	Social       []Social      `json:"social,omitempty"`
	Achievements []Achievement `json:"achievements,omitempty"`
	AverageAge   float32       `json:"averageAge"`
}

type Roaster []TeamPlayer

type TeamPlayer struct {
	Team       Team    `json:"-"`
	Nickname   string  `json:"nickname"`
	Status     string  `json:"status"`
	TimeOnTeam string  `json:"timeonteam"`
	MapsPlayed int     `json:"mapsplayed"`
	Rating     float32 `json:"rating"`
}
