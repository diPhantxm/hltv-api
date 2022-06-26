package main

import (
	"hltvapi/internal/api"
	"hltvapi/internal/controllers"
	"hltvapi/internal/controllers/storedImpl"
	"hltvapi/internal/controllers/storedImpl/store"
)

func main() {
	config := api.Config{
		Address: "localhost",
		Port:    ":8080",
	}

	dbConfig := store.Config{
		ConnectionString: "user=postgres password=3d1p4h1t5m dbname=hltvapi sslmode=disable",
		Driver:           "postgres",
	}

	store, err := store.NewStore(dbConfig)
	if err != nil {
		return
	}

	strategy := controllers.ControllersStrategy{
		Events:  storedImpl.NewEventsController(store),
		Matches: storedImpl.NewMatchesController(store),
		Players: storedImpl.NewPlayersController(store),
		Teams:   storedImpl.NewTeamsController(store),
	}

	api := api.NewApi(config, strategy)
	api.Run()
}
