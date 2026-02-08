// Package v1 provides HTTP handlers and DTOs for API version 1.
package v1

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"open-replays/internal/api/domain"
	"open-replays/internal/api/http/validator"
	"open-replays/internal/api/response"
	"open-replays/internal/api/usecase"
)

// VideosHandler is a handler for videos.
type VideosHandler struct {
	service   *usecase.VideosService
	validator *validator.UploadValidator
}

// NewVideosHandler creates a new VideosHandler.
func NewVideosHandler(
	service *usecase.VideosService,
	maxFileSize int64,
	allowedExts []string,
) *VideosHandler {
	return &VideosHandler{
		service:   service,
		validator: validator.NewUploadValidator(maxFileSize, allowedExts),
	}
}

// List lists all videos.
func (h *VideosHandler) List(c *gin.Context) {
	videos, err := h.service.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
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

	response.OK(c, summaries)
}

// GetByID gets a video by ID.
func (h *VideosHandler) GetByID(c *gin.Context) {
	var req GetVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	getByIDParams := usecase.GetByIDParams{ID: req.ID}
	video, err := h.service.GetByID(c.Request.Context(), getByIDParams)
	if err != nil {
		switch {
		case domain.IsNotFound(err):
			response.NotFound(c, "video")
		default:
			response.InternalError(c, err)
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

	response.OK(c, detail)
}

// Upload uploads a video.
func (h *VideosHandler) Upload(c *gin.Context) {
	var req UploadVideoRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	// Validate title
	if validationResult := h.validator.ValidateTitle(req.Title); !validationResult.IsValid() {
		response.Error(c, validationResult.ToError())
		return
	}

	// Validate video file
	if err := h.validator.ValidateFile(req.Video); err != nil {
		response.Error(c, err)
		return
	}

	uploadParams := usecase.UploadParams{
		File:        req.Video,
		Title:       req.Title,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Ext:         filepath.Ext(req.Video.Filename),
	}
	video, err := h.service.Upload(c.Request.Context(), uploadParams)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKWithMessage(c, UploadVideoResponse{ID: video.ID}, "Video uploaded successfully")
}

// Delete deletes a video.
func (h *VideosHandler) Delete(c *gin.Context) {
	var req DeleteVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	deleteParams := usecase.DeleteParams{ID: req.ID}
	if err := h.service.Delete(c.Request.Context(), deleteParams); err != nil {
		switch {
		case domain.IsNotFound(err):
			response.NotFound(c, "video")
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.OKWithMessage(c, DeleteVideoResponse(req), "Video deleted successfully")
}
