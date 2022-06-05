package models

type Player struct {
	Id        int        `json:"id"`
	Age       int        `json:"age"`
	Nickname  string     `json:"nickname"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Country   string     `json:"country"`
	Stats     Statistics `json:"stats"`
	Social    []Social   `json:"social,omitempty"`
}
