package models

type Achievement struct {
	Id        int
	Team      Team
	Player    Player
	Name      string
	Placement string
}
