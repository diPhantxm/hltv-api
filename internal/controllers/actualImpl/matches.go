package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
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

	parser := parsers.MatchParser{}

	event, err := parser.GetMatch(id)
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

	parser := parsers.MatchParser{}
	matches, err := parser.GetMatchesByDate(startDate)

	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, matches)
}

func (c MatchesController) GetAll(ctx *gin.Context) {
	parser := parsers.MatchParser{}

	matches, err := parser.GetMatches()
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, matches)
}
