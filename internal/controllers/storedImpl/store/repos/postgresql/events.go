package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
	"strings"
	"time"
)

type EventsRepo struct {
	db *sql.DB
}

func NewEventsRepo(db *sql.DB) EventsRepo {
	return EventsRepo{
		db: db,
	}
}

func (r EventsRepo) Get(expr func(event models.Event) bool) []models.Event {
	rows, err := r.db.Query(`SELECT id, name, startdate, enddate, prizepool, location, teams FROM events`)
	if err != nil {
		return nil
	}

	var allEvents []models.Event
	for rows.Next() {

		var event models.Event
		var teamsStr string
		var startDateTicks int64
		var endDateTicks int64

		if err := rows.Scan(&event.Id, &event.Name, &startDateTicks, &endDateTicks, &event.PrizePool, &event.Location, &teamsStr); err != nil {
			return nil
		}
		event.Teams = strings.Split(teamsStr, ",")
		event.StartDate = time.Unix(startDateTicks, 0)
		event.EndDate = time.Unix(endDateTicks, 0)

		allEvents = append(allEvents, event)
	}

	var result []models.Event
	for _, event := range allEvents {
		if expr(event) {
			result = append(result, event)
		}
	}

	return result
}

func (r EventsRepo) AddOrEdit(event models.Event) {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM events WHERE id=$1`, event.Id).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.edit(event)
	} else if count == 0 {
		r.add(event)
	}
}

func (r EventsRepo) add(event models.Event) {
	r.db.Exec(`INSERT INTO events (id, name, startdate, enddate, prizepool, location, teams) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.Id,
		event.Name,
		event.StartDate.Unix(),
		event.EndDate.Unix(),
		event.PrizePool,
		event.Location,
		strings.Join(event.Teams, ","))
}

func (r EventsRepo) edit(event models.Event) {
	r.db.Exec(`UPDATE events 
	SET name=$1, startDate=$2, endDate=$3, prizePool=$4, location=$5, teams=$6
	WHERE id=$7`,
		event.Name,
		event.StartDate.Unix(),
		event.EndDate.Unix(),
		event.PrizePool,
		event.Location,
		strings.Join(event.Teams, ","),
		event.Id)
}
