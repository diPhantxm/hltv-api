package models

type Player struct {
	Id           int           `json:"id"`
	Age          int           `json:"age"`
	Nickname     string        `json:"nickname"`
	Team         string        `json:"team"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Country      string        `json:"country"`
	Achievements []Achievement `json:"achievements"`
	Stats        Statistics    `json:"stats,omitempty"`
	Social       []Social      `json:"social,omitempty"`
}
