package repos

import (
	"hltvapi/internal/models"
)

type MatchesRepo interface {
	Get(expr func(models.Match) bool) []models.Match
	AddOrEdit(match models.Match)
}
