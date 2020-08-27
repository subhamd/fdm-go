package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type FilesController struct {
}

func (filesController FilesController) GetDownloadedFiles(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "downloaded files"})
}