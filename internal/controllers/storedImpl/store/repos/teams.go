package repos

import "hltvapi/internal/models"

type TeamsRepo interface {
	Get(expr func(models.Team) bool) []models.Team
	AddOrEdit(team models.Team)
}
