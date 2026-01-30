package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"open-replays/api/internal/api/domain"
	"open-replays/api/internal/api/usecase"
)

// VideosHandler is a handler for videos.
type VideosHandler struct {
	service *usecase.VideosService
}

// NewVideosHandler creates a new VideosHandler.
func NewVideosHandler(service *usecase.VideosService) *VideosHandler {
	return &VideosHandler{service: service}
}

// List lists all videos.
func (h *VideosHandler) List(c *gin.Context) {
	videos, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false, Error: "internal server error",
		})
		return
	}

	summaries := make([]VideoSummaryDTO, len(videos))
	for i, video := range videos {
		summaries[i] = VideoSummaryDTO{
			ID:           video.ID,
			Title:        video.Title,
			ThumbnailURL: video.ThumbnailURL,
			UploadedAt:   video.UploadedAt.Format(time.RFC3339),
			Duration:     video.Duration,
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true, Data: gin.H{"videos": summaries},
	})
}

// GetByID gets a video by ID.
func (h *VideosHandler) GetByID(c *gin.Context) {
	var req GetVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false, Error: err.Error(),
		})
		return
	}

	getByIDParams := usecase.GetByIDParams{ID: req.ID}
	video, err := h.service.GetByID(c.Request.Context(), getByIDParams)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrVideoNotFound) || errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false, Error: domain.ErrVideoNotFound.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false, Error: "internal error",
			})
		}
		return
	}

	detail := VideoDetailDTO{
		ID:           video.ID,
		Title:        video.Title,
		Description:  video.Description,
		ThumbnailURL: video.ThumbnailURL,
		VideoURL:     video.VideoURL,
		Duration:     video.Duration,
		Views:        video.Views,
		Author:       nil,
		Comments:     []CommentDTO{},
		UploadedAt:   video.UploadedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true, Data: detail,
	})
}

// Upload uploads a video.
func (h *VideosHandler) Upload(c *gin.Context) {
	var req UploadVideoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false, Error: err.Error(),
		})
		return
	}

	uploadParams := usecase.UploadParams{
		File:        req.Video,
		Title:       req.Title,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
	}
	video, err := h.service.Upload(c.Request.Context(), uploadParams)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidFileType):
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false, Error: domain.ErrInvalidFileType.Error(),
			})
		case errors.Is(err, domain.ErrFileTooLarge):
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false, Error: domain.ErrFileTooLarge.Error(),
			})
		case errors.Is(err, domain.ErrValidation):
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false, Error: domain.ErrValidation.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false, Error: "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true, Data: gin.H{"id": video.ID}, Message: "File uploaded successfully",
	})
}

// Delete deletes a video.
func (h *VideosHandler) Delete(c *gin.Context) {
	var req DeleteVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false, Error: err.Error(),
		})
		return
	}

	deleteParams := usecase.DeleteParams{ID: req.ID}
	err := h.service.Delete(c.Request.Context(), deleteParams)
	if err != nil {
		switch {
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false, Error: "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}
