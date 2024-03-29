package repos

import "hltvapi/internal/models"

type EventsRepo interface {
	Get(expr func(event models.Event) bool) []models.Event
	AddOrEdit(event models.Event)
}
