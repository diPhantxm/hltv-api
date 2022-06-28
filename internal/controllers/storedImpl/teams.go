package storedImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/controllers/storedImpl/store"
	"hltvapi/internal/controllers/storedImpl/taskScheduler"
	"hltvapi/internal/models"
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder/httpUrlBuilder"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TeamsController struct {
	store  *store.Store
	parser *parsers.TeamParser
}

func NewTeamsController(store *store.Store) *TeamsController {
	return &TeamsController{
		store:  store,
		parser: parsers.NewTeamParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
}

func (c TeamsController) GetAll(ctx *gin.Context) {
	teams := c.store.Teams.Get(func(team models.Team) bool {
		return true
	})

	ctx.JSON(http.StatusOK, teams)
}

func (c TeamsController) GetById(ctx *gin.Context) {
	// id: int
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No id parameter in request")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
		return
	}

	teams := c.store.Teams.Get(func(team models.Team) bool {
		return team.Id == id
	})

	if len(teams) != 1 {
		controllers.Error(ctx, http.StatusInternalServerError, "There was 0 or more than 1 entities in database")
		return
	}

	ctx.JSON(http.StatusOK, teams[0])
}

func (c TeamsController) Run() {
	for {
		ids, err := c.parser.GetAllTeamsIds()
		if err != nil {
			continue
		}

		for _, id := range ids {
			taskScheduler.GetStraightScheduler().Add(c.poll, id)
		}

		time.Sleep(5 * time.Minute)
	}
}

func (c TeamsController) poll(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	team, err := c.parser.GetTeam(id)
	if err != nil {
		return
	}

	c.store.Teams.AddOrEdit(*team)
}
