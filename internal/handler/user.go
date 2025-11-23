package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/utils"
	"gorm.io/gorm"
)

// UserHandler HTTP request handler for users
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Register a new user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user   body      dto.CreateUserRequest  true  "User Data"
// @Success      201    {object}  models.User
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		// Check for admin role blocking
		if err.Error() == "cannot register admin users via API" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		// Check for duplicate email
		if service.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200     {array}   models.User
// @Failure      500     {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get details of a specific user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	user, err := h.service.GetUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                true  "User ID"
// @Param        user   body      dto.UpdateUserRequest  true  "User Data"
// @Success      200    {object}  models.User
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// Check for duplicate email
		if service.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db save failed"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	err := h.service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterWithRole registers a user with a specific role (convenience method)
// Used for /register/owner and /register/consultant endpoints
func (h *UserHandler) RegisterWithRole(c *gin.Context, role models.UserRole) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override role with the specified one
	req.Role = role

	user, err := h.service.CreateUser(&req)
	if err != nil {
		// Check for duplicate email
		if service.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
