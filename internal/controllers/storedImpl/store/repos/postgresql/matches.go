package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
	"strings"
	"time"
)

type MatchesRepo struct {
	db *sql.DB
}

func (r MatchesRepo) Get(expr func(models.Match) bool) []models.Match {
	rows, err := r.db.Query(`SELECT id, teama, teamb, starttime, maps, viewers, playerofthematch FROM matches`)
	if err != nil {
		return nil
	}

	var allMatches []models.Match
	for rows.Next() {
		var match models.Match
		var startTime int64
		var maps string

		if err := rows.Scan(&match.Id, &match.TeamA, &match.TeamB, &startTime, &maps, &match.Viewers, &match.PlayerOfTheMatch); err != nil {
			return nil
		}

		match.StartTime = time.Unix(startTime, 0)
		match.Maps = strings.Split(maps, ",")

		allMatches = append(allMatches, match)
	}

	var result []models.Match
	for _, match := range allMatches {
		if expr(match) {
			result = append(result, match)
		}
	}

	return result
}

func (r MatchesRepo) AddOrEdit(match models.Match) {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM matches WHERE id=$1`, match.Id).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.Edit(match)
	} else if count == 0 {
		r.Add(match)
	}
}

func (r MatchesRepo) Add(match models.Match) {
	r.db.Exec(`INSERT INTO matches (id, teama, teamb, starttime, maps, viewers, playerofthematch) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		match.Id,
		match.TeamA,
		match.TeamB,
		match.StartTime.Unix(),
		strings.Join(match.Maps, ","),
		match.Viewers,
		match.PlayerOfTheMatch,
	)
}

func (r MatchesRepo) Edit(match models.Match) {
	r.db.Exec(`UPDATE matches 
	SET teama=$1, teamb=$2, starttime=$3, maps=$4, viewers=$5, playerofthematch=$6
	WHERE id=$7`,
		match.TeamA,
		match.TeamB,
		match.StartTime.Unix(),
		strings.Join(match.Maps, ","),
		match.Viewers,
		match.PlayerOfTheMatch,
		match.Id)
}
