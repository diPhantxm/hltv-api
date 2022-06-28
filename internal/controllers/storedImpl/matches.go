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

type MatchesController struct {
	store  *store.Store
	parser *parsers.MatchParser
}

func NewMatchesController(store *store.Store) *MatchesController {
	return &MatchesController{
		store:  store,
		parser: parsers.NewMatchParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
}

func (c MatchesController) GetById(ctx *gin.Context) {
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

	matches := c.store.Matches.Get(func(match models.Match) bool {
		return match.Id == id
	})

	if len(matches) != 1 {
		controllers.Error(ctx, http.StatusInternalServerError, "There was 0 or more than 1 entities in database")
		return
	}

	ctx.JSON(http.StatusOK, matches[0])
}

func (c MatchesController) GetByDate(ctx *gin.Context) {
	// date: time.Time
	dateStr, ok := ctx.Params.Get("date")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No date parameter in request")
		return
	}

	layout := `02-01-2006`
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
		return
	}

	matches := c.store.Matches.Get(func(match models.Match) bool {
		return match.StartTime.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour))
	})

	ctx.JSON(http.StatusOK, matches)
}

func (c MatchesController) GetAll(ctx *gin.Context) {
	matches := c.store.Matches.Get(func(match models.Match) bool {
		return true
	})

	ctx.JSON(http.StatusOK, matches)
}

func (c MatchesController) Run() {
	for {
		ids, err := c.parser.GetUpcomingMatchesIds()
		if err != nil {
			continue
		}

		for _, id := range ids {
			taskScheduler.GetStraightScheduler().Add(c.poll, id)
		}

		time.Sleep(5 * time.Minute)
	}
}

func (c MatchesController) poll(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	match, err := c.parser.GetMatch(id)
	if err != nil {
		return
	}

	c.store.Matches.AddOrEdit(*match)
}
