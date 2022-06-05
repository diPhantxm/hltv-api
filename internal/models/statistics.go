package models

type Statistics struct {
	Rating            float32 `json:"rating"`   // Rating 2.0
	KillsPerRound     float32 `json:"kr"`       // Kills per round
	Headshots         float32 `json:"headshot"` // Headshot percentage
	MapsPlayed        int     `json:"mapsPlayed"`
	DeathsPerRound    float32 `json:"dr"` // Deaths per round
	RoundsContributed float32 `json:"roundsContributed"`
}
