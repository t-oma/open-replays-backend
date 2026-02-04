package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"open-replays/internal/api/domain"
	"open-replays/internal/api/usecase"
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
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

	c.JSON(http.StatusOK, SuccessResponse[[]VideoSummaryDTO]{
		Data: summaries,
	})
}

// GetByID gets a video by ID.
func (h *VideosHandler) GetByID(c *gin.Context) {
	var req GetVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	getByIDParams := usecase.GetByIDParams{ID: req.ID}
	video, err := h.service.GetByID(c.Request.Context(), getByIDParams)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrVideoNotFound) || errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: domain.ErrVideoNotFound.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "internal error",
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

	c.JSON(http.StatusOK, SuccessResponse[VideoDetailDTO]{
		Data: detail,
	})
}

// Upload uploads a video.
func (h *VideosHandler) Upload(c *gin.Context) {
	var req UploadVideoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
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
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: domain.ErrInvalidFileType.Error(),
			})
		case errors.Is(err, domain.ErrFileTooLarge):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: domain.ErrFileTooLarge.Error(),
			})
		case errors.Is(err, domain.ErrValidation):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: domain.ErrValidation.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse[UploadVideoResponse]{
		Data: UploadVideoResponse{ID: video.ID}, Message: "File uploaded successfully",
	})
}

// Delete deletes a video.
func (h *VideosHandler) Delete(c *gin.Context) {
	var req DeleteVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	deleteParams := usecase.DeleteParams{ID: req.ID}
	err := h.service.Delete(c.Request.Context(), deleteParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse[DeleteVideoResponse]{
		Data:    DeleteVideoResponse{ID: req.ID},
		Message: "Video deleted successfully",
	})
}
