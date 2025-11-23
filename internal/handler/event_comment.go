package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/middleware"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/storage"
	"github.com/you/pawtrack/internal/utils"
	"gorm.io/gorm"
)

type EventCommentHandler struct {
	service service.EventCommentService
	storage storage.FileStorage
}

func NewEventCommentHandler(service service.EventCommentService, storage storage.FileStorage) *EventCommentHandler {
	return &EventCommentHandler{
		service: service,
		storage: storage,
	}
}

// CreateComment godoc
// @Summary      Create event comment
// @Description  Create a new comment on an event (Owner or Consultant with access)
// @Tags         event-comments
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Event ID"
// @Param        comment  body      dto.CreateCommentRequest  true  "Comment Data"
// @Success      201      {object}  models.EventComment
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /events/{id}/comments [post]
func (h *EventCommentHandler) CreateComment(c *gin.Context) {
	eventID := uint(utils.Atoi(c.Param("id")))

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if this is multipart or JSON
	contentType := c.GetHeader("Content-Type")
	var req dto.CreateCommentRequest

	if contentType == "application/json" {
		// Legacy JSON support
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
			return
		}

		// Parse JSON data from form field
		dataStr := c.PostForm("data")
		if dataStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing data field"})
			return
		}

		if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in data field"})
			return
		}

		// Handle file upload if present
		file, header, err := c.Request.FormFile("file")
		if err == nil {
			defer file.Close()

			// Validate file
			if err := utils.ValidateFile(file, header); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Upload file
			fileURL, err := h.storage.Upload(file, header.Filename, header.Header.Get("Content-Type"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
				return
			}

			req.AttachmentURL = &fileURL
		}
	}

	// Override event_id from URL
	req.EventID = eventID

	comment, err := h.service.CreateComment(&req, userID, role)
	if err != nil {
		if err.Error() == "not authorized" || err.Error() == "event has no associated dog" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// ListComments godoc
// @Summary      List event comments
// @Description  Get all comments for an event
// @Tags         event-comments
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  dto.CommentListResponse
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /events/{id}/comments [get]
func (h *EventCommentHandler) ListComments(c *gin.Context) {
	eventID := uint(utils.Atoi(c.Param("id")))

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.service.ListComments(eventID, userID, role)
	if err != nil {
		if err.Error() == "not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comments"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetComment godoc
// @Summary      Get event comment
// @Description  Get comment by ID
// @Tags         event-comments
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  dto.CommentResponse
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /event-comments/{id} [get]
func (h *EventCommentHandler) GetComment(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	comment, err := h.service.GetComment(id, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		if err.Error() == "not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment godoc
// @Summary      Update event comment
// @Description  Update comment by ID (only author or admin)
// @Tags         event-comments
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Comment ID"
// @Param        comment  body      dto.UpdateCommentRequest  true  "Comment Data"
// @Success      200      {object}  models.EventComment
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /event-comments/{id} [put]
func (h *EventCommentHandler) UpdateComment(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.UpdateComment(id, &req, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		if err.Error() == "only comment author can update" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary      Delete event comment
// @Description  Delete comment by ID (only author or admin)
// @Tags         event-comments
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      204  {object}  nil
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /event-comments/{id} [delete]
func (h *EventCommentHandler) DeleteComment(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.service.DeleteComment(id, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		if err.Error() == "only comment author can delete" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	c.Status(http.StatusNoContent)
}
