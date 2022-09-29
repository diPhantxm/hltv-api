package store

import (
	"database/sql"
	"fmt"
	"hltvapi/internal/controllers/storedImpl/store/repos"
)

type Store struct {
	db      *sql.DB
	Events  repos.EventsRepo
	Matches repos.MatchesRepo
	Teams   repos.TeamsRepo
	Players repos.PlayersRepo
}

func NewStore(config Config) (*Store, error) {
	db, err := sql.Open(config.Driver, config.ConnectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	factory := NewRepoFactory(config, db)

	s := &Store{
		Events:  factory.CreateEventsRepo(),
		Matches: factory.CreateMatchesRepo(),
		Teams:   factory.CreateTeamsRepo(),
		Players: factory.CreatePlayersRepo(),
	}

	s.db = db

	return s, nil
}

func (s *Store) Close() {
	s.db.Close()
}
