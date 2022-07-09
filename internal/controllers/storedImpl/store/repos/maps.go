package repos

import "hltvapi/internal/models"

type MapsRepo interface {
	Get(expr func(matchMap models.Map) bool) []models.Map
	GetByMatchId(id int) []models.Map
	AddOrEdit(matchMap models.Map)
	Add(matchMap models.Map)
	Edit(matchMap models.Map)
}
