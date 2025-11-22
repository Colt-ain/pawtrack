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

// DogHandler HTTP request handler for dogs
type DogHandler struct {
	service service.DogService
}

// NewDogHandler creates a new dog handler
func NewDogHandler(service service.DogService) *DogHandler {
	return &DogHandler{service: service}
}

// CreateDog godoc
// @Summary      Create a new dog
// @Description  Create a new dog (requires authentication and owner role)
// @Tags         dogs
// @Accept       json
// @Produce      json
// @Param        dog    body      dto.CreateDogRequest  true  "Dog Data"
// @Success      201    {object}  models.Dog
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /dogs [post]
func (h *DogHandler) CreateDog(c *gin.Context) {
	var req dto.CreateDogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	dog, err := h.service.CreateDog(&req, userID.(uint))
	if err != nil {
		// Check for date parsing error
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, use RFC3339"})
		return
	}

	c.JSON(http.StatusCreated, dog)
}

// ListDogs godoc
// @Summary      List dogs
// @Description  Get a list of dogs
// @Tags         dogs
// @Produce      json
// @Success      200     {array}   models.Dog
// @Failure      500     {object}  map[string]string
// @Router       /dogs [get]
func (h *DogHandler) ListDogs(c *gin.Context) {
	dogs, err := h.service.ListDogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, dogs)
}

// GetDog godoc
// @Summary      Get dog by ID
// @Description  Get details of a specific dog
// @Tags         dogs
// @Produce      json
// @Param        id   path      int  true  "Dog ID"
// @Success      200  {object}  models.Dog
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /dogs/{id} [get]
func (h *DogHandler) GetDog(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	dog, err := h.service.GetDog(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, dog)
}

// UpdateDog godoc
// @Summary      Update dog
// @Description  Update a dog by ID
// @Tags         dogs
// @Accept       json
// @Produce      json
// @Param        id     path      int               true  "Dog ID"
// @Param        dog    body      dto.UpdateDogRequest  true  "Dog Data"
// @Success      200    {object}  models.Dog
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dogs/{id} [put]
func (h *DogHandler) UpdateDog(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	var req dto.UpdateDogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dog, err := h.service.UpdateDog(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, use RFC3339"})
		return
	}

	c.JSON(http.StatusOK, dog)
}

// DeleteDog godoc
// @Summary      Delete dog
// @Description  Delete a dog by ID
// @Tags         dogs
// @Produce      json
// @Param        id   path      int  true  "Dog ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /dogs/{id} [delete]
func (h *DogHandler) DeleteDog(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	err := h.service.DeleteDog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
