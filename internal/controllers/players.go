package controllers

import "github.com/gin-gonic/gin"

type PlayersController interface {
	Controller
	GetById(ctx *gin.Context)
	GetByNickname(ctx *gin.Context)
}
