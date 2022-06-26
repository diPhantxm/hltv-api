package repos

import (
	"hltvapi/internal/models"
)

type MatchesRepo interface {
	Get(expr func(models.Match) bool) []models.Match
	AddOrEdit(match models.Match)
	Add(match models.Match)
	Edit(match models.Match)
}
