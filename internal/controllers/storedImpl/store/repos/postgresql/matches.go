package postgresql

import (
	"database/sql"
	"fmt"
	"hltvapi/internal/controllers/storedImpl/store/repos"
	"hltvapi/internal/models"
	"time"
)

type MatchesRepo struct {
	db       *sql.DB
	mapsRepo repos.MapsRepo
}

func NewMatchesRepo(db *sql.DB) MatchesRepo {
	return MatchesRepo{
		db:       db,
		mapsRepo: NewMapsRepo(db),
	}
}

func (r MatchesRepo) Get(expr func(models.Match) bool) []models.Match {
	rows, err := r.db.Query(`SELECT id, teama, teamb, starttime, viewers, playerofthematch FROM matches`)
	if err != nil {
		return nil
	}

	var allMatches []models.Match
	for rows.Next() {
		var match models.Match
		var startTime int64

		if err := rows.Scan(&match.Id, &match.TeamA, &match.TeamB, &startTime, &match.Viewers, &match.PlayerOfTheMatch); err != nil {
			return nil
		}

		match.StartTime = time.Unix(startTime, 0)

		allMatches = append(allMatches, match)
	}

	var result []models.Match
	for _, match := range allMatches {
		if expr(match) {
			match.Maps = r.mapsRepo.GetByMatchId(match.Id)

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
	_, err := r.db.Exec(`INSERT INTO matches (id, teama, teamb, starttime, viewers, playerofthematch, isover) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		match.Id,
		match.TeamA,
		match.TeamB,
		match.StartTime.Unix(),
		match.Viewers,
		match.PlayerOfTheMatch,
		match.IsOver,
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, matchMap := range match.Maps {
		matchMap.MatchId = match.Id
		r.mapsRepo.Add(matchMap)
	}
}

func (r MatchesRepo) Edit(match models.Match) {
	r.db.Exec(`UPDATE matches 
	SET teama=$1, teamb=$2, starttime=$3, viewers=$5, playerofthematch=$6
	WHERE id=$7`,
		match.TeamA,
		match.TeamB,
		match.StartTime.Unix(),
		match.Viewers,
		match.PlayerOfTheMatch,
		match.Id)

	for _, matchMap := range match.Maps {
		r.mapsRepo.Edit(matchMap)
	}
}
