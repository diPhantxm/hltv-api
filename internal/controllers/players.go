package controllers

import "github.com/gin-gonic/gin"

type PlayersController interface {
	Controller
	GetById(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}
