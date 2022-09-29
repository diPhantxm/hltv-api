package models

type Social struct {
	Team   Team
	Player Player
	Name   string `json:"name"`
	Link   string `json:"link"`
}
