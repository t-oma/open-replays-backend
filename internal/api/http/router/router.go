package router

import (
	"github.com/gin-gonic/gin"
	v1 "open-replays/api/internal/api/http/v1"
	"open-replays/api/internal/api/usecase"
)

func Register(r *gin.Engine, todos *usecase.VideosService) {
	api := r.Group("/api")
	v1g := api.Group("/v1")

	h := v1.NewVideosHandler(todos)
	v1g.GET("/videos", h.List)
	v1g.POST("/videos/upload", h.Upload)
	v1g.DELETE("/videos/:filename", h.Delete)
	v1g.GET("/videos/:filename/watch", h.Watch)
}
