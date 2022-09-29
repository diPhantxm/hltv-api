package postgresql

import (
	"database/sql"
	"hltvapi/internal/controllers/storedImpl/store/repos"
)

type RepoFactory struct {
	db *sql.DB
}

func NewFactory(db *sql.DB) *RepoFactory {
	return &RepoFactory{
		db: db,
	}
}

func (f *RepoFactory) CreateEventsRepo() repos.EventsRepo {
	return NewEventsRepo(f.db)
}

func (f *RepoFactory) CreateMatchesRepo() repos.MatchesRepo {
	return NewMatchesRepo(f.db)
}

func (f *RepoFactory) CreatePlayersRepo() repos.PlayersRepo {
	return NewPlayersRepo(f.db)
}

func (f *RepoFactory) CreateTeamsRepo() repos.TeamsRepo {
	return NewTeamsRepo(f.db)
}
