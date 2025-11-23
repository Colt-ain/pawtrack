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

// EventHandler HTTP request handler for events
type EventHandler struct {
	service service.EventService
	storage storage.FileStorage
}

// NewEventHandler creates a new event handler
func NewEventHandler(service service.EventService, storage storage.FileStorage) *EventHandler {
	return &EventHandler{
		service: service,
		storage: storage,
	}
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Create a new pet event with optional file attachment
// @Tags         events
// @Accept       multipart/form-data
// @Produce      json
// @Param        data  formData  string  true  "Event Data (JSON)"
// @Param        file  formData  file    false "File attachment"
// @Success      201   {object}  models.Event
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// Check if this is multipart or JSON
	contentType := c.GetHeader("Content-Type")
	var req dto.CreateEventRequest

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

	event, err := h.service.CreateEvent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db create failed"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// ListEvents godoc
// @Summary      List events with filtering
// @Description  Get paginated events with filters and sorting
// @Tags         events
// @Produce      json
// @Param        from_date    query     string  false  "From date (YYYY-MM-DD)"
// @Param        to_date      query     string  false  "To date (YYYY-MM-DD)"
// @Param        types        query     string  false  "Event types (comma-separated)"
// @Param        search       query     string  false  "Search in notes and dog names"
// @Param        dog_name     query     string  false  "Filter by dog name"
// @Param        page         query     int     false  "Page number" default(1)
// @Param        page_size    query     int     false  "Page size" default(20)
// @Param        sort_by      query     string  false  "Sort by field" Enums(created_at, type)
// @Param        sort_order   query     string  false  "Sort order" Enums(asc, desc)
// @Success      200          {object}  dto.EventListResponse
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	var filters dto.EventFilterParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole, err := middleware.GetUserRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filters.UserID = userID
	filters.UserRole = userRole

	response, err := h.service.ListEvents(&filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetEvent godoc
// @Summary      Get event by ID
// @Description  Get details of a specific event
// @Tags         events
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	event, err := h.service.GetEvent(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Delete event
// @Description  Delete an event by ID
// @Tags         events
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	err := h.service.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
