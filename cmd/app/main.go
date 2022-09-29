package main

import (
	"hltvapi/internal/api"
	"hltvapi/internal/controllers"
	"hltvapi/internal/controllers/storedImpl"
	"hltvapi/internal/controllers/storedImpl/store"
	"log"
)

func main() {
	config := api.Config{
		Address: "0.0.0.0",
		Port:    ":5000",
	}

	dbConfig := store.Config{
		ConnectionString: "host=hltv-db port=5432 user=postgres password=3d1p4h1t5m dbname=hltvapi sslmode=disable",
		Driver:           "postgres",
	}

	store, err := store.NewStore(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
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
