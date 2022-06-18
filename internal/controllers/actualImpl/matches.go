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

type MatchesController struct {
}

func (c MatchesController) Run() {

}

func (c MatchesController) GetById(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No id parameter in request")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
	}

	urlBuilder := urlBuilder.NewUrlBuilder()
	urlBuilder.Match()
	urlBuilder.AddId(id)

	event, err := parsers.GetMatch(urlBuilder.String())
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c MatchesController) GetByDate(ctx *gin.Context) {
	startDate, ok := ctx.Params.Get("date")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "Date was not found in request")
		return
	}

	url := urlBuilder.NewUrlBuilder()
	url.Results()
	url.AddParam("startDate", startDate)
	url.AddParam("endDate", startDate)

	matchesIds, err := parsers.GetMatchesIdsByDate(url.String())
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	matches := make([]models.Match, len(matchesIds))
	for i, id := range matchesIds {
		url = urlBuilder.NewUrlBuilder()
		url.Match()
		url.AddId(id)

		match, err := parsers.GetMatch(url.String())
		if err != nil {
			continue
		}
		matches[i] = *match
	}

	ctx.JSON(http.StatusOK, matches)
}

func (c MatchesController) GetAll(ctx *gin.Context) {
	url := urlBuilder.NewUrlBuilder()
	url.Match()
	matchesListLink := url.String()

	upcomingMatchesIds, err := parsers.GetUpcomingMatchesIds(matchesListLink)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	matches := make([]models.Match, len(upcomingMatchesIds))
	for i, id := range upcomingMatchesIds {
		matchUrl := urlBuilder.NewUrlBuilder()
		matchUrl.Match()
		matchUrl.AddId(id)

		match, err := parsers.GetMatch(matchUrl.String())
		if err != nil {
			controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}

		matches[i] = *match
	}

	ctx.JSON(http.StatusOK, matches)
}
