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

type PlayersController struct {
	store  *store.Store
	parser *parsers.PlayerParser
}

func NewPlayersController(store *store.Store) *PlayersController {
	return &PlayersController{
		store:  store,
		parser: parsers.NewPlayerParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
}

func (c PlayersController) GetAll(ctx *gin.Context) {
	players := c.store.Players.Get(func(player models.Player) bool {
		return true
	})

	ctx.JSON(http.StatusOK, players)
}

func (c PlayersController) GetById(ctx *gin.Context) {
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

	players := c.store.Players.Get(func(player models.Player) bool {
		return player.Id == id
	})

	if len(players) != 1 {
		controllers.Error(ctx, http.StatusInternalServerError, "There was 0 or more than 1 entities in database")
		return
	}

	ctx.JSON(http.StatusOK, players[0])
}

func (c PlayersController) Run() {
	for {
		ids, err := c.parser.GetAllPlayersIds()
		if err != nil {
			continue
		}

		for _, id := range ids {
			taskScheduler.GetStraightScheduler().Add(c.poll, id)
		}

		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func (c PlayersController) poll(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	player, err := c.parser.GetPlayer(id)
	if err != nil {
		return
	}

	c.store.Players.AddOrEdit(*player)
}
