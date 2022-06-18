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

	urlBuilder := urlBuilder.NewUrlBuilder()
	urlBuilder.Player()
	urlBuilder.AddId(id)

	event, err := parsers.GetPlayer(urlBuilder.String())
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c PlayersController) GetAll(ctx *gin.Context) {
	url := urlBuilder.NewUrlBuilder()
	url.PlayersStats()
	playersStatsList := url.String()

	playersIds, err := parsers.GetAllPlayersIds(playersStatsList)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	players := make([]models.Player, len(playersIds))
	for i, id := range playersIds {
		playerUrl := urlBuilder.NewUrlBuilder()
		playerUrl.Player()
		playerUrl.AddId(id)

		player, err := parsers.GetPlayer(playerUrl.String())
		if err != nil {
			controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}

		players[i] = *player
	}

	ctx.JSON(http.StatusOK, players)
}
