package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamsController struct {
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

	parser := parsers.TeamParser{}
	event, err := parser.GetTeam(id)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c TeamsController) GetAll(ctx *gin.Context) {
	parser := parsers.TeamParser{}
	teams, err := parser.GetTeams()

	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, teams)
}
