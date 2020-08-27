package main

import (
	"github.com/gin-gonic/gin"
	"github.com/subhamd/fdm-go/controllers"
	"log"
	"net/http"
)

func main() {
	app := new(controllers.AppController)

	router := gin.Default()
	router = setupRoutes(router, *app)
	runApp(router, *app)
}

func setupRoutes(router *gin.Engine, app controllers.AppController) *gin.Engine {
	router.GET("/health", app.HealthCheck)

	downloadsController := new(controllers.DownloadsController)
	filesController := new(controllers.FilesController)

	v1 := router.Group("/v1")

	downloadsGroup := v1.Group("/downloads")
	downloadsGroup.POST("", downloadsController.PostDownloads)
	downloadsGroup.GET("/:id", downloadsController.GetStatusByDownloadId)

	filesGroup := v1.Group("/files")
	filesGroup.GET("", filesController.GetDownloadedFiles)

	return router
}

func runApp(router *gin.Engine, app controllers.AppController) {
	listenPort := "8081"

	server := &http.Server{
		Addr:    "0.0.0.0:" + listenPort,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println(err)
		//panic(err)
	}
}
