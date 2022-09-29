package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type TeamsRepo struct {
	db               *sql.DB
	socialsRepo      SocialsRepo
	achievementsRepo AchievementsRepo
	roastersRepo     RoastersRepo
}

func NewTeamsRepo(db *sql.DB) TeamsRepo {
	return TeamsRepo{
		db:               db,
		socialsRepo:      NewSocialsRepo(db),
		achievementsRepo: NewAchievementsRepo(db),
		roastersRepo:     NewRoastersRepo(db),
	}
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
			team.Social = r.socialsRepo.GetByTeamId(team.Id)
			team.Achievements = r.achievementsRepo.GetByTeamId(team.Id)
			team.Roaster = r.roastersRepo.GetByTeamId(team.Id)

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
		r.edit(team)
	} else if count == 0 {
		r.add(team)
	}

	for _, social := range team.Social {
		social.Team = team
		r.socialsRepo.AddOrEdit(social)
	}

	for _, achievement := range team.Achievements {
		achievement.Team = team
		r.achievementsRepo.AddOrEdit(achievement)
	}

	for i := range team.Roaster {
		team.Roaster[i].Team = team
	}

	r.roastersRepo.AddOrEdit(team.Roaster)
}

func (r TeamsRepo) add(team models.Team) {
	r.db.Exec(`INSERT INTO teams (id, ranking, weeksintop30, averageage, name, country) VALUES ($1, $2, $3, $4, $5, $6)`,
		team.Id,
		team.Ranking,
		team.WeeksInTop30,
		team.AverageAge,
		team.Name,
		team.Country,
	)
}

func (r TeamsRepo) edit(team models.Team) {
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
