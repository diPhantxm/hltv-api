package store

import (
	"database/sql"
	"hltvapi/internal/controllers/storedImpl/store/repos"
	"hltvapi/internal/controllers/storedImpl/store/repos/postgresql"

	_ "github.com/lib/pq"
)

type RepoFactory interface {
	CreateEventsRepo() repos.EventsRepo
	CreateMatchesRepo() repos.MatchesRepo
	CreatePlayersRepo() repos.PlayersRepo
	CreateTeamsRepo() repos.TeamsRepo
}

func NewRepoFactory(config Config, db *sql.DB) RepoFactory {
	switch config.Driver {
	case "postgres":
		return postgresql.NewFactory(db)
	}

	return nil
}
