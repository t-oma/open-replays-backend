// Package v1 provides HTTP handlers and DTOs for API version 1.
package v1

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"open-replays/internal/api/domain"
	"open-replays/internal/api/response"
	"open-replays/internal/api/usecase"
	"open-replays/internal/api/validation"
)

// VideosHandler is a handler for videos.
type VideosHandler struct {
	service *usecase.VideosService
	rules   validation.VideoRules
}

// NewVideosHandler creates a new VideosHandler.
func NewVideosHandler(
	service *usecase.VideosService,
	rules validation.VideoRules,
) *VideosHandler {
	return &VideosHandler{
		service: service,
		rules:   rules,
	}
}

// List lists all videos.
func (h *VideosHandler) List(c *gin.Context) {
	var req ListVideosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	var errors []response.ValidationError

	errors = append(errors, validation.Int64("page", int64(req.Page)).
		Min(1).
		Collect()...)

	errors = append(errors, validation.Int64("pageSize", int64(req.PageSize)).
		Min(1).
		Max(100).
		Collect()...)

	if len(errors) > 0 {
		response.ValidationFailed(c, errors)
		return
	}

	listParams := usecase.ListParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	videos, err := h.service.List(c.Request.Context(), listParams)
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

	totalItems, err := h.service.Count(c.Request.Context())
	if err != nil {
		switch {
		case domain.IsNotFound(err):
			response.NotFound(c, "videos")
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.OK(c, response.PaginatedData[VideoSummaryDTO]{
		Items:      summaries,
		Pagination: response.NewPagination(req.Page, req.PageSize, totalItems),
	})
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

	// Validate using fluent API
	var errors []response.ValidationError

	errors = append(errors, validation.String("title", req.Title).
		Required().
		MinLength(h.rules.MinTitleLength).
		MaxLength(h.rules.MaxTitleLength).
		Collect()...)

	errors = append(errors, validation.String("description", req.Description).
		Optional().
		MaxLength(h.rules.MaxDescriptionLength).
		Collect()...)

	errors = append(errors, validation.File("video", req.Video).
		Required().
		MinSize(h.rules.MinVideoSize).
		MaxSize(h.rules.MaxVideoSize).
		AllowedExts(h.rules.AllowedVideoExts).
		Collect()...)

	errors = append(errors, validation.File("thumbnail", req.Thumbnail).
		Optional().
		MinSize(h.rules.MinThumbnailSize).
		MaxSize(h.rules.MaxThumbnailSize).
		AllowedExts(h.rules.AllowedThumbnailExts).
		Collect()...)

	if len(errors) > 0 {
		response.ValidationFailed(c, errors)
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
