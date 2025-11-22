package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/utils"
	"gorm.io/gorm"
)

// EventHandler HTTP request handler for events
type EventHandler struct {
	service service.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Create a new pet event
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event  body      dto.CreateEventRequest  true  "Event Data"
// @Success      201    {object}  models.Event
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db create failed"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// ListEvents godoc
// @Summary      List events
// @Description  Get a list of events with optional filtering
// @Tags         events
// @Produce      json
// @Param        limit   query     int     false  "Limit"  default(50)
// @Param        offset  query     int     false  "Offset" default(0)
// @Param        type    query     string  false  "Filter by type"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]string
// @Router       /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	limit := utils.ClampAtoi(c.DefaultQuery("limit", "50"), 1, 200)
	offset := utils.Max(utils.Atoi(c.DefaultQuery("offset", "0")), 0)
	eventType := c.Query("type")

	events, err := h.service.ListEvents(limit, offset, eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": events, "limit": limit, "offset": offset})
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
