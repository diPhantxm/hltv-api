package controllers

import "github.com/gin-gonic/gin"

type EventsController interface {
	Controller
	GetById(ctx *gin.Context)
	GetByName(ctx *gin.Context)
	GetByPrizePool(ctx *gin.Context)
	GetByCountry(ctx *gin.Context)
}
