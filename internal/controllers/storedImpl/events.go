package storedImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/controllers/storedImpl/store"
	"hltvapi/internal/controllers/storedImpl/taskScheduler"
	"hltvapi/internal/models"
	"hltvapi/internal/parsers"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EventsController struct {
	store *store.Store
}

func NewEventsController(store *store.Store) *EventsController {
	return &EventsController{
		store: store,
	}
}

func (c EventsController) GetById(ctx *gin.Context) {
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

	events := c.store.Events.Get(func(event models.Event) bool {
		return event.Id == id
	})

	if len(events) != 1 {
		controllers.Error(ctx, http.StatusInternalServerError, "There was 0 or more than 1 entities in the database")
		return
	}

	ctx.JSON(http.StatusOK, events[0])
}

func (c EventsController) GetAll(ctx *gin.Context) {
	events := c.store.Events.Get(func(event models.Event) bool {
		return true
	})

	ctx.JSON(http.StatusOK, events)
}

func (c EventsController) Run() {
	p := parsers.EventParser{}

	for {
		ids, err := p.GetUpcomingEventsIds()
		if err != nil {
			continue
		}

		for _, id := range ids {
			taskScheduler.GetStraightScheduler().Add(c.poll, id)
		}

		time.Sleep(5 * time.Minute)
	}
}

func (c EventsController) poll(params ...interface{}) {
	id := params[0].([]interface{})[0].(int)

	p := parsers.EventParser{}

	event, err := p.GetEvent(id)
	if err != nil {
		return
	}

	c.store.Events.AddOrEdit(*event)
}
