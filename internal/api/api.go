package api

import (
	"hltvapi/internal/controllers"

	"github.com/gin-gonic/gin"
)

type Api struct {
	teamsController   controllers.TeamsController
	playersController controllers.PlayersController
	matchesController controllers.MatchesController
	eventsController  controllers.EventsController

	handler *gin.Engine
	config  Config
}

func NewApi(config Config, strategy controllers.ControllersStrategy) *Api {
	return &Api{
		teamsController:   strategy.Teams,
		playersController: strategy.Players,
		matchesController: strategy.Matches,
		eventsController:  strategy.Events,
		handler:           gin.New(),
		config:            config,
	}
}

func (a Api) Run() {
	a.teamsController.Run()
	a.playersController.Run()
	a.matchesController.Run()
	a.eventsController.Run()

	a.configureRoutes()

	a.handler.Run(a.config.Port)
}

func (a Api) configureRoutes() {
	v1 := a.handler.Group("/v1")
	{
		teams := v1.Group("teams/")
		{
			teams.GET("/", a.teamsController.GetAll)
			teams.GET("/:id", a.teamsController.GetById)
		}

		players := v1.Group("/players")
		{
			players.GET("/", a.playersController.GetAll)
			players.GET("/:id", a.playersController.GetById)
		}

		matches := v1.Group("/matches")
		{
			matches.GET("/", a.matchesController.GetAll)
			matches.GET("/:id", a.matchesController.GetById)
			matches.GET("/date/:date", a.matchesController.GetByDate)
		}

		events := v1.Group("/events")
		{
			events.GET("/", a.eventsController.GetAll)
			events.GET("/:id", a.eventsController.GetById)
		}
	}
}
