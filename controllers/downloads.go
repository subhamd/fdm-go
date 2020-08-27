package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/subhamd/fdm-go/entities"
	"github.com/subhamd/fdm-go/services"
	"log"
	"net/http"
)

type DownloadsController struct {
}

func (downloadsController DownloadsController) PostDownloads(context *gin.Context) {
	var request entities.PostDownloadsRequestObject

	err := context.BindJSON(&request)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"internal_code": 1001, "message": "Invalid request body"})
		return
	}

	if err = request.TYPE.IsValid(); err != nil {
		log.Println("Unknown download type")
		context.JSON(
			http.StatusBadRequest,
			gin.H{"internal_code": 1002, "message": err.Error()})
		return
	}

	downloadStatus, downloadId := services.DownloadFiles(request.TYPE, request.URLS)
	if downloadStatus == entities.DownloadStatusFailed {
		context.JSON(http.StatusInternalServerError, gin.H{"internal_code": 1003, "message": "Oops! Some went wrong. Unable to download files at this moment"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": downloadId})
}

func (downloadsController DownloadsController) GetStatusByDownloadId(context *gin.Context) {
	downloadId := context.Param("id")

	downloadEntity, err := services.GetDownloadStatus(downloadId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"internal_code": 1004, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, downloadEntity)
}