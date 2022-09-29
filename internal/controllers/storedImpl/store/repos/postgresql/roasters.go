package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type RoastersRepo struct {
	db *sql.DB
}

func NewRoastersRepo(db *sql.DB) RoastersRepo {
	return RoastersRepo{
		db: db,
	}
}

func (r RoastersRepo) GetByTeamId(id int) models.Roaster {
	rows, err := r.db.Query(`SELECT nickname, status, timeonteam, mapsplayed, rating
							FROM roasters
							WHERE teamid=$1`, id)

	if err != nil {
		return nil
	}

	var roaster models.Roaster

	for rows.Next() {
		var player models.TeamPlayer

		if err := rows.Scan(&player.Nickname, &player.Status, &player.TimeOnTeam, &player.MapsPlayed, &player.Rating); err != nil {
			continue
		}

		roaster = append(roaster, player)
	}

	return roaster
}

func (r RoastersRepo) AddOrEdit(roaster models.Roaster) {
	for _, player := range roaster {
		var count int
		if err := r.db.QueryRow(`SELECT COUNT(*) FROM roasters WHERE nickname=$1`, player.Nickname).Scan(&count); err != nil {
			return
		}

		if count == 1 {
			r.edit(player)
		} else if count == 0 {
			r.add(player)
		}
	}
}

func (r RoastersRepo) add(player models.TeamPlayer) {
	r.db.Exec(`INSERT INTO roasters (nickname, status, timeonteam, mapsplayed, rating, teamid) VALUES ($1, $2, $3, $4, $5, $6)`,
		player.Nickname,
		player.Status,
		player.TimeOnTeam,
		player.MapsPlayed,
		player.Rating,
		player.Team.Id,
	)
}

func (r RoastersRepo) edit(player models.TeamPlayer) {
	r.db.Exec(`UPDATE roasters SET status=$1, timeonteam=$2, mapsplayed=$3, rating=$4
	WHERE teamid=$5 AND nickname=$6`,
		player.Status,
		player.TimeOnTeam,
		player.MapsPlayed,
		player.Rating,
		player.Team.Id,
		player.Nickname,
	)
}
