package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/middleware"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/utils"
	"gorm.io/gorm"
)

type ConsultantNoteHandler struct {
	service service.ConsultantNoteService
}

func NewConsultantNoteHandler(service service.ConsultantNoteService) *ConsultantNoteHandler {
	return &ConsultantNoteHandler{service: service}
}

// CreateNote godoc
// @Summary      Create consultant note
// @Description  Create a new note about a dog (Consultant only, requires access to dog)
// @Tags         consultant-notes
// @Accept       json
// @Produce      json
// @Param        note  body      dto.CreateNoteRequest  true  "Note Data"
// @Success      201   {object}  models.ConsultantNote
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /consultant-notes [post]
func (h *ConsultantNoteHandler) CreateNote(c *gin.Context) {
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

	// Only consultants can create notes
	if role != models.RoleConsultant && role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only consultants can create notes"})
		return
	}

	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.CreateNote(&req, userID)
	if err != nil {
		if err.Error() == "consultant does not have access to this dog" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNote godoc
// @Summary      Get consultant note
// @Description  Get note by ID
// @Tags         consultant-notes
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Success      200  {object}  dto.NoteResponse
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /consultant-notes/{id} [get]
func (h *ConsultantNoteHandler) GetNote(c *gin.Context) {
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

	note, err := h.service.GetNote(id, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// UpdateNote godoc
// @Summary      Update consultant note
// @Description  Update note by ID
// @Tags         consultant-notes
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "Note ID"
// @Param        note  body      dto.UpdateNoteRequest   true  "Note Data"
// @Success      200   {object}  models.ConsultantNote
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /consultant-notes/{id} [put]
func (h *ConsultantNoteHandler) UpdateNote(c *gin.Context) {
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

	var req dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.UpdateNote(id, &req, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote godoc
// @Summary      Delete consultant note
// @Description  Delete note by ID
// @Tags         consultant-notes
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Success      204  {object}  nil
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /consultant-notes/{id} [delete]
func (h *ConsultantNoteHandler) DeleteNote(c *gin.Context) {
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

	err = h.service.DeleteNote(id, userID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListNotes godoc
// @Summary      List consultant notes
// @Description  List notes with filtering and sorting
// @Tags         consultant-notes
// @Produce      json
// @Param        search    query     string  false  "Search in title and content"
// @Param        dog_id    query     int     false  "Filter by dog ID"
// @Param        owner_id  query     int     false  "Filter by owner ID"
// @Param        from_date query     string  false  "Filter from date (RFC3339)"
// @Param        to_date   query     string  false  "Filter to date (RFC3339)"
// @Param        sort_by   query     string  false  "Sort by: created_at, updated_at, dog_name, owner_name" default(created_at)
// @Param        order     query     string  false  "Order: asc, desc" default(desc)
// @Param        page      query     int     false  "Page number" default(1)
// @Param        page_size query     int     false  "Page size" default(20)
// @Success      200       {object}  dto.NoteListResponse
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Security     BearerAuth
// @Router       /consultant-notes [get]
func (h *ConsultantNoteHandler) ListNotes(c *gin.Context) {
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

	var filters dto.NoteFilterParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListNotes(&filters, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notes"})
		return
	}

	c.JSON(http.StatusOK, result)
}
