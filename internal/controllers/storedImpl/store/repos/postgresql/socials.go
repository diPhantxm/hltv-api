package postgresql

import (
	"database/sql"
	"hltvapi/internal/models"
)

type SocialsRepo struct {
	db *sql.DB
}

func NewSocialsRepo(db *sql.DB) SocialsRepo {
	return SocialsRepo{
		db: db,
	}
}

func (r SocialsRepo) GetByTeamId(id int) []models.Social {
	rows, err := r.db.Query(`SELECT name, link FROM socials WHERE teamid=$1`, id)
	if err != nil {
		return nil
	}

	var socials []models.Social
	for rows.Next() {
		var social models.Social

		if err := rows.Scan(&social.Name, &social.Link); err != nil {
			return nil
		}

		socials = append(socials, social)
	}

	return socials
}

func (r SocialsRepo) GetByPlayerId(id int) []models.Social {
	rows, err := r.db.Query(`SELECT name, link FROM socials WHERE playerid=$1`, id)
	if err != nil {
		return nil
	}

	var socials []models.Social
	for rows.Next() {
		var social models.Social

		if err := rows.Scan(&social.Name, &social.Link); err != nil {
			return nil
		}

		socials = append(socials, social)
	}

	return socials
}

func (r SocialsRepo) AddOrEdit(social models.Social) {
	var count int
	if err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM socials 
		WHERE teamid=$1 OR 
		playerid=$2 AND 
		name=$3`,
		social.Team.Id, social.Player.Id, social.Name).Scan(&count); err != nil {
		return
	}

	if count == 1 {
		r.edit(social)
	} else if count == 0 {
		r.add(social)
	}
}

func (r SocialsRepo) add(social models.Social) {
	_, err := r.db.Exec(`INSERT INTO socials (name, link, teamid, playerid) VALUES ($1, $2, $3, $4)`,
		social.Name,
		social.Link,
		social.Team.Id,
		social.Player.Id,
	)

	if err != nil {

	}
}

func (r SocialsRepo) edit(social models.Social) {

}
