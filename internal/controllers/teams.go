package controllers

import "github.com/gin-gonic/gin"

type TeamsController interface {
	Controller
	GetById(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}
