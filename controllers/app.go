package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppController struct {
}

func (appController AppController) HealthCheck(context  *gin.Context) {
	context.JSON(http.StatusOK, gin.H{})
}
