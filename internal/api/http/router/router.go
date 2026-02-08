// Package router provides HTTP route registration.
package router

import (
	"github.com/gin-gonic/gin"

	v1 "open-replays/internal/api/http/v1"
)

// Register registers all HTTP routes.
func Register(
	r *gin.Engine,
	videosHandler *v1.VideosHandler,
) {
	api := r.Group("/api")
	v1g := api.Group("/v1")

	v1g.GET("/videos", videosHandler.List)
	v1g.GET("/videos/:id", videosHandler.GetByID)
	v1g.POST("/videos/upload", videosHandler.Upload)

	r.Static("/media", "./uploads")
}
