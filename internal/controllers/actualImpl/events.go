package actualImpl

import (
	"hltvapi/internal/controllers"
	"hltvapi/internal/parsers"
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

	parser := parsers.EventParser{}

	event, err := parser.GetEvent(id)
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, event)
}

func (c EventsController) GetAll(ctx *gin.Context) {
	parser := parsers.EventParser{}

	events, err := parser.GetEvents()
	if err != nil {
		controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, events)
}
