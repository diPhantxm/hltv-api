package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
	"hltvapi/internal/urlBuilder/httpUrlBuilder"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventsController struct {
	parser *parsers.EventParser
}

func NewEventsController() *EventsController {
	return &EventsController{
		parser: parsers.NewEventParser(httpUrlBuilder.NewHttpUrlBuilder()),
	}
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

	event, err := c.parser.GetEvent(id)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c EventsController) GetAll(ctx *gin.Context) {
	events, err := c.parser.GetEvents()
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, events)
}
