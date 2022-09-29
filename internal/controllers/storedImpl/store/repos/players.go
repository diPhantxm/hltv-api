package repos

import (
	"hltvapi/internal/models"
)

type PlayersRepo interface {
	Get(expr func(models.Player) bool) []models.Player
	AddOrEdit(player models.Player)
}
