package controllers

type ControllersStrategy struct {
	Events  EventsController
	Matches MatchesController
	Teams   TeamsController
	Players PlayersController
}
