package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/middleware"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/utils"
)

type ConsultantHandler struct {
	service service.ConsultantService
}

func NewConsultantHandler(service service.ConsultantService) *ConsultantHandler {
	return &ConsultantHandler{service: service}
}

// UpdateProfile godoc
// @Summary      Update consultant profile
// @Description  Update profile details (Consultant only)
// @Tags         consultants
// @Accept       json
// @Produce      json
// @Param        profile  body      dto.UpdateProfileRequest  true  "Profile Data"
// @Success      200      {object}  models.ConsultantProfile
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /consultants/profile [put]
func (h *ConsultantHandler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetProfile godoc
// @Summary      Get consultant profile
// @Description  Get profile by ID
// @Tags         consultants
// @Produce      json
// @Param        id   path      int  true  "Consultant ID"
// @Success      200  {object}  dto.ConsultantProfileResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /consultants/{id} [get]
func (h *ConsultantHandler) GetProfile(c *gin.Context) {
	id := uint(utils.Atoi(c.Param("id")))

	profile, err := h.service.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// SearchConsultants godoc
// @Summary      Search consultants
// @Description  Search consultants by name, services, breeds, location
// @Tags         consultants
// @Produce      json
// @Param        query     query     string  false  "Search query"
// @Param        services  query     string  false  "Filter by services"
// @Param        breeds    query     string  false  "Filter by breeds"
// @Param        location  query     string  false  "Filter by location"
// @Param        page      query     int     false  "Page number" default(1)
// @Param        page_size query     int     false  "Page size" default(20)
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /consultants [get]
func (h *ConsultantHandler) SearchConsultants(c *gin.Context) {
	var req dto.ConsultantSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profiles, count, err := h.service.SearchConsultants(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        profiles,
		"total_count": count,
		"page":        req.Page,
		"page_size":   req.PageSize,
	})
}

// InviteConsultant godoc
// @Summary      Invite consultant
// @Description  Invite a consultant to manage a dog (Owner only)
// @Tags         consultants
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Consultant ID"
// @Param        request  body      dto.CreateInviteRequest  true  "Invite Data"
// @Success      201      {object}  dto.InviteResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /consultants/{id}/invite [post]
func (h *ConsultantHandler) InviteConsultant(c *gin.Context) {
	consultantID := uint(utils.Atoi(c.Param("id")))

	ownerID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invite, err := h.service.InviteConsultant(ownerID, consultantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invite"})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// AcceptInvite godoc
// @Summary      Accept invite
// @Description  Accept an invitation to manage a dog (Consultant only)
// @Tags         invites
// @Produce      json
// @Param        token    query     string  true  "Invite Token"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /invites/accept [post]
func (h *ConsultantHandler) AcceptInvite(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}

	consultantID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.service.AcceptInvite(token, consultantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invite accepted"})
}
