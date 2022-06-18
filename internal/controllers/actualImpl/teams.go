package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/models"
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder"
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

	urlBuilder := urlBuilder.NewUrlBuilder()
	urlBuilder.Team()
	urlBuilder.AddId(id)

	event, err := parsers.GetTeam(urlBuilder.String())
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c TeamsController) GetAll(ctx *gin.Context) {
	url := urlBuilder.NewUrlBuilder()
	url.TeamsStats()
	teamsStatsList := url.String()

	teamsIds, err := parsers.GetAllTeamsIds(teamsStatsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	teams := make([]models.Team, len(teamsIds))
	for i, id := range teamsIds {
		teamUrl := urlBuilder.NewUrlBuilder()
		teamUrl.Team()
		teamUrl.AddId(id)

		team, err := parsers.GetTeam(teamUrl.String())
		if err != nil {
			controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}

		teams[i] = *team
	}

	ctx.JSON(http.StatusOK, teams)
}
