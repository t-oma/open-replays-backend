package v1

import (
	"database/sql"
	"errors"
	"fmt"
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
		c.JSON(500, APIResponse{
			Success: false,
			Code:    500,
			Error:   err.Error(),
		})
		return
	}

	summaries := make([]VideoSummaryDTO, len(videos))
	for i, video := range videos {
		summaries[i] = VideoSummaryDTO{
			ID:           video.ID,
			Title:        video.Title,
			ThumbnailURL: "/thumbnails/" + video.Filename + ".jpg",
			UploadedAt:   video.UploadedAt.Format(time.RFC3339),
			Duration:     video.Duration,
		}
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    gin.H{"videos": summaries},
	})
}

// GetByID gets a video by ID.
func (h *VideosHandler) GetByID(c *gin.Context) {
	var req GetVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Code:    400,
			Error:   err.Error(),
		})
		return
	}

	getByIDParams := usecase.GetByIDParams{ID: req.ID}
	video, err := h.service.GetByID(c.Request.Context(), getByIDParams)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound,
				APIResponse{
					Success: false,
					Code:    http.StatusNotFound,
					Error:   domain.ErrVideoNotFound.Error(),
				})
			return
		}
		fmt.Println(err)

		c.JSON(500, APIResponse{Success: false, Error: "internal error"})
		return
	}

	detail := VideoDetailDTO{
		ID:           video.ID,
		Title:        video.Title,
		Description:  video.Description,
		ThumbnailURL: "/thumbnails/" + video.Thumbnail,
		VideoURL:     "/" + video.ID,
		Duration:     video.Duration,
		Views:        video.Views,
		// Author:       &AuthorDTO{},
		// Comments:     []CommentDTO{},
		UploadedAt: video.UploadedAt.Format(time.RFC3339),
	}

	c.JSON(200, APIResponse{Success: true, Data: detail})
}

// Upload uploads a video.
func (h *VideosHandler) Upload(c *gin.Context) {
	var req UploadVideoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Code:    400,
			Error:   err.Error(),
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
			c.JSON(400, APIResponse{Success: false, Code: 400, Error: err.Error()})
		case errors.Is(err, domain.ErrFileTooLarge):
			c.JSON(400, APIResponse{Success: false, Code: 400, Error: err.Error()})
		case errors.Is(err, domain.ErrValidation):
			c.JSON(400, APIResponse{Success: false, Code: 400, Error: err.Error()})
		default:
			c.JSON(500, APIResponse{Success: false, Code: 500, Error: "Internal server error"})
		}
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    gin.H{"id": video.ID},
		Message: "File uploaded successfully",
	})
}

// Delete deletes a video.
func (h *VideosHandler) Delete(c *gin.Context) {
	var req DeleteVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Code:    400,
			Error:   err.Error(),
		})
		return
	}

	deleteParams := usecase.DeleteParams{ID: req.ID}
	err := h.service.Delete(c.Request.Context(), deleteParams)
	if err != nil {
		switch {
		default:
			c.JSON(500, APIResponse{
				Success: false,
				Code:    500,
				Error:   err.Error(),
			})
			return
		}
	}

	c.JSON(200, APIResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}

// Watch watch a video
func (h *VideosHandler) Watch(c *gin.Context) {
	var req WatchVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Code:    400,
			Error:   err.Error(),
		})
		return
	}

	watchParams := usecase.WatchParams{ID: req.ID}
	reader, video, err := h.service.Watch(c.Request.Context(), watchParams)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound), errors.Is(err, domain.ErrVideoNotFound), errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Code:    http.StatusNotFound,
				Error:   err.Error(),
			})
			return
		default:
			c.JSON(500, APIResponse{
				Success: false,
				Code:    500,
				Error:   err.Error(),
			})
			return
		}
	}
	defer reader.Close()

	contentType := getContentType(video.Extension)
	filename := fmt.Sprintf("%s%s", video.Filename, video.Extension)

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	c.Header("Content-Type", contentType)

	c.DataFromReader(http.StatusOK, -1, contentType, reader, nil)
}

func getContentType(ext string) string {
	switch ext {
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	case "mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}
