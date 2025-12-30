package v1

import (
	"errors"

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

	c.JSON(200, APIResponse{
		Success: true,
		Data:    gin.H{"videos": videos},
	})
}

// GetByFilename gets a video by filename.
func (h *VideosHandler) GetByFilename(c *gin.Context) {
	var req GetVideoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Code:    400,
			Error:   err.Error(),
		})
		return
	}

	getByFilenameParams := usecase.GetByFilenameParams{Filename: req.Filename}
	video, err := h.service.GetByFilename(c.Request.Context(), getByFilenameParams)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Code:    500,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    gin.H{"video": video},
	})
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
		Data:    gin.H{"filename": video.Filename + "." + video.Extension},
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

	deleteParams := usecase.DeleteParams{Filename: req.Filename}
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

	watchParams := usecase.WatchParams{Filename: req.Filename}
	file, err := h.service.Watch(c.Request.Context(), watchParams)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			c.JSON(404, APIResponse{
				Success: false,
				Code:    404,
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

	c.File(file)
}
