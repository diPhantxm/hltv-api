package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayersController struct {
}

func (c PlayersController) Run() {

}

func (c PlayersController) GetById(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No id parameter in request")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
	}

	parser := parsers.PlayerParser{}

	event, err := parser.GetPlayer(id)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c PlayersController) GetAll(ctx *gin.Context) {
	parser := parsers.PlayerParser{}

	players, err := parser.GetPlayers()
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, players)
}
