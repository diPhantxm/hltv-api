package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type TeamsRepo struct {
	db *sql.DB
}

func (r TeamsRepo) Get(expr func(models.Team) bool) []models.Team {
	rows, err := r.db.Query(`SELECT id, ranking, weeksintop30, averageage, name, country
	FROM teams`)
	if err != nil {
		return nil
	}

	var allTeams []models.Team
	for rows.Next() {
		var team models.Team

		if err := rows.Scan(&team.Id, &team.Ranking, &team.WeeksInTop30, &team.AverageAge, &team.Name, &team.Country); err != nil {
			return nil
		}

		allTeams = append(allTeams, team)
	}

	var result []models.Team
	for _, team := range allTeams {
		if expr(team) {
			result = append(result, team)
		}
	}

	return result
}

func (r TeamsRepo) AddOrEdit(team models.Team) {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM teams WHERE id=$1`, team.Id).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.Edit(team)
	} else if count == 0 {
		r.Add(team)
	}
}

func (r TeamsRepo) Add(team models.Team) {
	r.db.Exec(`INSERT INTO teams (id, ranking, weeksintop30, averageage, name, country) VALUES ($1, $2, $3, $4, $5, $6)`,
		team.Id,
		team.Ranking,
		team.WeeksInTop30,
		team.AverageAge,
		team.Name,
		team.Country,
	)
}

func (r TeamsRepo) Edit(team models.Team) {
	r.db.Exec(`UPDATE teams SET ranking=$1, weeksintop30=$2, averageage=$3, name=$4, country=$5
	WHERE id=$6`,
		team.Ranking,
		team.WeeksInTop30,
		team.AverageAge,
		team.Name,
		team.Country,
		team.Id,
	)
}
