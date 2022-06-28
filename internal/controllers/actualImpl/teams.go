package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder/httpUrlBuilder"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamsController struct {
	parser *parsers.TeamParser
}

func NewTeamsController() *TeamsController {
	return &TeamsController{
		parser: parsers.NewTeamParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
}

func (c TeamsController) Run() {

}

func (c TeamsController) GetById(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No id parameter in request")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
	}

	event, err := c.parser.GetTeam(id)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c TeamsController) GetAll(ctx *gin.Context) {
	teams, err := c.parser.GetTeams()

	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, teams)
}
