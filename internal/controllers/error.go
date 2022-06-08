package controllers

import "github.com/gin-gonic/gin"

func Error(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"Error": message,
	})
}
