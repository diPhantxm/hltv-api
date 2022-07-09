package postgresql

import (
	"database/sql"
	"fmt"
	"hltvapi/internal/models"
)

type MapsRepo struct {
	db *sql.DB
}

func NewMapsRepo(db *sql.DB) MapsRepo {
	return MapsRepo{
		db: db,
	}
}

func (r MapsRepo) Get(expr func(matchMap models.Map) bool) []models.Map {
	rows, err := r.db.Query(`SELECT id, matchid, name, teamascore, teambscore FROM maps`)
	if err != nil {
		return nil
	}

	var allMaps []models.Map
	for rows.Next() {
		var matchMap models.Map

		if err := rows.Scan(&matchMap.Id, &matchMap.MatchId, &matchMap.Name, &matchMap.TeamAScore, &matchMap.TeamBScore); err != nil {
			return nil
		}

		allMaps = append(allMaps, matchMap)
	}

	var result []models.Map
	for _, matchMap := range allMaps {
		if expr(matchMap) {
			result = append(result, matchMap)
		}
	}

	return result
}

func (r MapsRepo) GetByMatchId(id int) []models.Map {
	rows, err := r.db.Query(`SELECT matchid, name, teamascore, teambscore FROM maps WHERE id=$1`, id)
	if err != nil {
		return nil
	}

	var maps []models.Map
	for rows.Next() {
		var matchMap models.Map
		matchMap.Id = id

		if err := rows.Scan(&matchMap.MatchId, &matchMap.Name, &matchMap.TeamAScore, &matchMap.TeamBScore); err != nil {
			return nil
		}

		maps = append(maps, matchMap)
	}

	return maps
}

func (r MapsRepo) AddOrEdit(matchMap models.Map) {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM maps WHERE id=$1`, matchMap.Id).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.Edit(matchMap)
	} else if count == 0 {
		r.Add(matchMap)
	}
}

func (r MapsRepo) Add(matchMap models.Map) {
	_, err := r.db.Exec(`INSERT INTO maps (matchid, name, teamascore, teambscore) VALUES ($1, $2, $3, $4)`,
		matchMap.MatchId,
		matchMap.Name,
		matchMap.TeamAScore,
		matchMap.TeamBScore)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (r MapsRepo) Edit(matchMap models.Map) {
	r.db.Exec(`UPDATE maps SET matchid=$1, name=$2, teamascore=$3, teambscore=$4 WHERE id=$5`,
		matchMap.MatchId,
		matchMap.Name,
		matchMap.TeamAScore,
		matchMap.TeamBScore,
		matchMap.Id)
}
