package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type PlayersRepo struct {
	db *sql.DB
}

func (r PlayersRepo) Get(expr func(models.Player) bool) []models.Player {
	rows, err := r.db.Query(`SELECT players.id, age, nickname, firstname, lastname, country, 
	rating, killsperround, headshots, mapsplayed, deathsperround, roundscontributed
	FROM players FULL JOIN stats ON players.id = stats.id`)
	if err != nil {
		return nil
	}

	var allPlayers []models.Player
	for rows.Next() {
		var player models.Player

		if err := rows.Scan(&player.Id, &player.Age, &player.Nickname, &player.FirstName, &player.LastName, &player.Country,
			&player.Stats.Rating, &player.Stats.KillsPerRound, &player.Stats.Headshots, &player.Stats.MapsPlayed,
			&player.Stats.DeathsPerRound, &player.Stats.RoundsContributed); err != nil {
			return nil
		}

		allPlayers = append(allPlayers, player)
	}

	var result []models.Player
	for _, player := range allPlayers {
		if expr(player) {
			result = append(result, player)
		}
	}

	return result
}

func (r PlayersRepo) AddOrEdit(player models.Player) {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM player WHERE id=$1`, player.Id).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.Edit(player)
	} else if count == 0 {
		r.Add(player)
	}
}

func (r PlayersRepo) Add(player models.Player) {
	r.db.Exec(`INSERT INTO players (id, age, nickname, firstname, lastname, country) VALUES ($1, $2, $3, $4, $5, $6)`,
		player.Id,
		player.Age,
		player.Nickname,
		player.FirstName,
		player.LastName,
		player.Country,
	)

	r.db.Exec(`INSERT INTO stats (id, rating, killsperround, headshots, mapsplayed, deathsperround, roundscontributed) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		player.Id,
		player.Stats.Rating,
		player.Stats.KillsPerRound,
		player.Stats.Headshots,
		player.Stats.MapsPlayed,
		player.Stats.DeathsPerRound,
		player.Stats.RoundsContributed,
	)
}

func (r PlayersRepo) Edit(player models.Player) {
	r.db.QueryRow(`UPDATE players SET age=$1, nickname=$2, firstname=$3, lastname=$4, country=$5
	WHERE id=$6`,
		player.Age,
		player.Nickname,
		player.FirstName,
		player.LastName,
		player.Country,
		player.Id,
	)

	r.db.QueryRow(`UPDATE stats SET rating=$1, killsperround=$2, headshots=$3, mapsplayed=$4, deathsperround=$5, roundscontributed=$6
	WHERE id=$7`,
		player.Stats.Rating,
		player.Stats.KillsPerRound,
		player.Stats.Headshots,
		player.Stats.MapsPlayed,
		player.Stats.DeathsPerRound,
		player.Stats.RoundsContributed,
		player.Id,
	)
}
