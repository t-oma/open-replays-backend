package router

import (
	"github.com/gin-gonic/gin"
	v1 "open-replays/api/internal/api/http/v1"
	"open-replays/api/internal/api/usecase"
)

// Register registers all HTTP routes.
func Register(r *gin.Engine, videosService *usecase.VideosService) {
	api := r.Group("/api")
	v1g := api.Group("/v1")

	h := v1.NewVideosHandler(videosService)
	v1g.GET("/videos", h.List)
	v1g.GET("/videos/:id", h.GetByID)
	v1g.POST("/videos/upload", h.Upload)

	r.Static("/media", "./uploads")
}
