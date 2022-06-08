package controllers

import "github.com/gin-gonic/gin"

type TeamsController interface {
	Controller
	GetById(ctx *gin.Context)
	GetByName(ctx *gin.Context)
	GetByCountry(ctx *gin.Context)
}
