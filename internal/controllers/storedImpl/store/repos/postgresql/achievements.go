package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type AchievementsRepo struct {
	db *sql.DB
}

func NewAchievementsRepo(db *sql.DB) AchievementsRepo {
	return AchievementsRepo{
		db: db,
	}
}

func (r AchievementsRepo) GetByPlayerId(id int) []models.Achievement {
	rows, err := r.db.Query(`SELECT name, placement FROM achievements WHERE playerid=$1`, id)
	if err != nil {
		return nil
	}

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement

		if err := rows.Scan(&achievement.Name, &achievement.Placement); err != nil {
			return nil
		}

		achievements = append(achievements, achievement)
	}

	return achievements
}

func (r AchievementsRepo) GetByTeamId(id int) []models.Achievement {
	rows, err := r.db.Query(`SELECT name, placement FROM achievements WHERE teamid=$1`, id)
	if err != nil {
		return nil
	}

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement

		if err := rows.Scan(&achievement.Name, &achievement.Placement); err != nil {
			return nil
		}

		achievements = append(achievements, achievement)
	}

	return achievements
}

func (r AchievementsRepo) AddOrEdit(achievements ...models.Achievement) {
	for _, achievement := range achievements {
		var count int
		if err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM achievements 
		WHERE teamid=$1 OR 
		playerid=$2 AND 
		name=$3`,
			achievement.Team.Id, achievement.Player.Id, achievement.Name).Scan(&count); err != nil {
			return
		}

		if count == 1 {
			r.edit(achievement)
		} else if count == 0 {
			r.add(achievement)
		}
	}
}

func (r AchievementsRepo) add(achievement models.Achievement) {
	_, err := r.db.Exec(`INSERT INTO achievements (name, placement, teamid, playerid) VALUES ($1, $2, $3, $4)`,
		achievement.Name,
		achievement.Placement,
		achievement.Team.Id,
		achievement.Player.Id,
	)

	if err != nil {

	}
}

func (r AchievementsRepo) edit(achievement models.Achievement) {

}
