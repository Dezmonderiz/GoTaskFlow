package handler

import (
	"net/http"

	"GoTaskFlow/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskService *service.TaskService) *gin.Engine {
	router := gin.Default()
	taskHandler := NewTaskHandler(taskService)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/styles.css", "./web/styles.css")
	router.StaticFile("/app.js", "./web/app.js")

	api := router.Group("/api")
	{
		api.GET("/tasks", taskHandler.List)
		api.POST("/tasks", taskHandler.Create)
		api.GET("/tasks/:id", taskHandler.GetByID)
		api.PATCH("/tasks/:id/status", taskHandler.UpdateStatus)
		api.DELETE("/tasks/:id", taskHandler.Delete)
		api.GET("/stats", taskHandler.Stats)
	}

	return router
}
