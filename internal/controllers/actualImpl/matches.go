package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder/httpUrlBuilder"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchesController struct {
	parser *parsers.MatchParser
}

func NewMatchesController() *MatchesController {
	return &MatchesController{
		parser: parsers.NewMatchParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
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

	event, err := c.parser.GetMatch(id)
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

	matches, err := c.parser.GetMatchesByDate(startDate)

	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, matches)
}

func (c MatchesController) GetAll(ctx *gin.Context) {
	matches, err := c.parser.GetMatches()
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, matches)
}
