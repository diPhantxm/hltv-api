package main

import (
	"hltvapi/internal/api"
	"hltvapi/internal/controllers"
	"hltvapi/internal/controllers/actualImpl"
)

func main() {
	config := api.Config{
		Address: "localhost",
		Port:    ":8080",
	}

	strategy := controllers.ControllersStrategy{
		Events:  actualImpl.EventsController{},
		Matches: actualImpl.MatchesController{},
		Players: actualImpl.PlayersController{},
		Teams:   actualImpl.TeamsController{},
	}

	api := api.NewApi(config, strategy)
	api.Run()
}
