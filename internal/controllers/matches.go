package controllers

import "github.com/gin-gonic/gin"

type MatchesController interface {
	Controller
	GetById(ctx *gin.Context)
	GetByDate(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}
