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

type EventsController struct {
}

func (c EventsController) Run() {

}

func (c EventsController) GetById(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		controllers.Error(ctx, http.StatusBadRequest, "No id parameter in request")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		controllers.Error(ctx, http.StatusBadRequest, "Id cannot be converted to int")
	}

	urlBuilder := urlBuilder.NewUrlBuilder()
	urlBuilder.Event()
	urlBuilder.AddId(id)

	event, err := parsers.GetEvent(urlBuilder.String())
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c EventsController) GetAll(ctx *gin.Context) {
	url := urlBuilder.NewUrlBuilder()
	url.Event()
	eventsListLink := url.String()

	upcomingEventsIds, err := parsers.GetUpcomingEventsIds(eventsListLink)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	events := make([]models.Event, len(upcomingEventsIds))
	for i, id := range upcomingEventsIds {
		eventUrl := urlBuilder.NewUrlBuilder()
		eventUrl.Event()
		eventUrl.AddId(id)

		event, err := parsers.GetEvent(eventUrl.String())
		if err != nil {
			controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}

		events[i] = *event
	}

	ctx.JSON(http.StatusOK, events)
}
