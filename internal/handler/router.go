package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/middleware"
	"github.com/you/pawtrack/internal/service"
	"github.com/you/pawtrack/internal/permissions"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all application routes
func SetupRouter(
	eventHandler *EventHandler,
	dogHandler *DogHandler,
	userHandler *UserHandler,
	authHandler *AuthHandler,
	healthHandler *HealthHandler,
	consultantHandler *ConsultantHandler,
	consultantNoteHandler *ConsultantNoteHandler,
	eventCommentHandler *EventCommentHandler,
	authService service.AuthService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Serve uploaded files (for local storage)
	router.Static("/uploads", "./uploads")

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/manual/swagger.json", "/srv/docs/swagger.json")

	// API v1
	api := router.Group("/api/v1")
	{
		// Public auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register/owner", func(c *gin.Context) {
				userHandler.RegisterWithRole(c, "owner")
			})
			auth.POST("/register/consultant", func(c *gin.Context) {
				userHandler.RegisterWithRole(c, "consultant")
			})
		}

		// Public users endpoint (for backward compatibility)
		api.POST("/users", userHandler.CreateUser)

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Events - require authentication
			protected.POST("/events", middleware.RequireAnyPermission(permissions.EVENTS_CREATE_OWN, permissions.EVENTS_CREATE_ASSIGNED, permissions.EVENTS_CREATE_ALL), eventHandler.CreateEvent)
			protected.GET("/events", middleware.RequireAnyPermission(permissions.EVENTS_VIEW_OWN, permissions.EVENTS_VIEW_ASSIGNED, permissions.EVENTS_VIEW_ALL), eventHandler.ListEvents)
			protected.GET("/events/:id", middleware.RequireAnyPermission(permissions.EVENTS_VIEW_OWN, permissions.EVENTS_VIEW_ASSIGNED, permissions.EVENTS_VIEW_ALL), eventHandler.GetEvent)
			protected.DELETE("/events/:id", middleware.RequireAnyPermission(permissions.EVENTS_DELETE_OWN, permissions.EVENTS_DELETE_ALL), eventHandler.DeleteEvent)

			// Dogs - create requires owner role, others just authentication
			protected.POST("/dogs", middleware.RequirePermission(permissions.DOGS_CREATE), dogHandler.CreateDog)
			protected.GET("/dogs", middleware.RequireAnyPermission(permissions.DOGS_VIEW_OWN, permissions.DOGS_VIEW_ASSIGNED, permissions.DOGS_VIEW_ALL), dogHandler.ListDogs)
			protected.GET("/dogs/:id", middleware.RequireAnyPermission(permissions.DOGS_VIEW_OWN, permissions.DOGS_VIEW_ASSIGNED, permissions.DOGS_VIEW_ALL), dogHandler.GetDog)
			protected.PUT("/dogs/:id", middleware.RequireAnyPermission(permissions.DOGS_UPDATE_OWN, permissions.DOGS_UPDATE_ALL), dogHandler.UpdateDog)
			protected.DELETE("/dogs/:id", middleware.RequireAnyPermission(permissions.DOGS_DELETE_OWN, permissions.DOGS_DELETE_ALL), dogHandler.DeleteDog)

			// Users - require authentication
			protected.GET("/users", middleware.RequireAnyPermission(permissions.USERS_VIEW_OWN, permissions.USERS_VIEW_ALL), userHandler.ListUsers)
			protected.GET("/users/:id", middleware.RequireAnyPermission(permissions.USERS_VIEW_OWN, permissions.USERS_VIEW_ALL), userHandler.GetUser)
			protected.PUT("/users/:id", middleware.RequireAnyPermission(permissions.USERS_UPDATE_OWN, permissions.USERS_UPDATE_ALL), userHandler.UpdateUser)
			protected.DELETE("/users/:id", middleware.RequirePermission(permissions.USERS_DELETE_ALL), userHandler.DeleteUser)

			// Consultants - require authentication
			protected.PUT("/consultants/profile", middleware.RequirePermission(permissions.CONSULTANTS_PROFILE_UPDATE), consultantHandler.UpdateProfile)
			protected.GET("/consultants", middleware.RequirePermission(permissions.CONSULTANTS_SEARCH), consultantHandler.SearchConsultants)
			protected.GET("/consultants/:id", middleware.RequirePermission(permissions.CONSULTANTS_SEARCH), consultantHandler.GetProfile)
			protected.POST("/consultants/:id/invite", middleware.RequirePermission(permissions.CONSULTANTS_INVITE), consultantHandler.InviteConsultant)

			// Invites - require authentication
			protected.POST("/invites/accept", middleware.RequirePermission(permissions.CONSULTANTS_INVITES_ACCEPT), consultantHandler.AcceptInvite)

			// Consultant Notes - require authentication
			protected.POST("/consultant-notes", middleware.RequirePermission(permissions.CONSULTANT_NOTES_CREATE), consultantNoteHandler.CreateNote)
			protected.GET("/consultant-notes", middleware.RequireAnyPermission(permissions.CONSULTANT_NOTES_VIEW_OWN, permissions.CONSULTANT_NOTES_VIEW_ALL), consultantNoteHandler.ListNotes)
			protected.GET("/consultant-notes/:id", middleware.RequireAnyPermission(permissions.CONSULTANT_NOTES_VIEW_OWN, permissions.CONSULTANT_NOTES_VIEW_ALL), consultantNoteHandler.GetNote)
			protected.PUT("/consultant-notes/:id", middleware.RequirePermission(permissions.CONSULTANT_NOTES_UPDATE_OWN), consultantNoteHandler.UpdateNote)
			protected.DELETE("/consultant-notes/:id", middleware.RequirePermission(permissions.CONSULTANT_NOTES_DELETE_OWN), consultantNoteHandler.DeleteNote)

			// Event Comments - require authentication
			protected.POST("/events/:id/comments", middleware.RequireAnyPermission(permissions.EVENT_COMMENTS_CREATE_OWN, permissions.EVENT_COMMENTS_CREATE_ASSIGNED), eventCommentHandler.CreateComment)
			protected.GET("/events/:id/comments", middleware.RequireAnyPermission(permissions.EVENT_COMMENTS_VIEW_OWN, permissions.EVENT_COMMENTS_VIEW_ASSIGNED), eventCommentHandler.ListComments)
			protected.GET("/event-comments/:id", middleware.RequireAnyPermission(permissions.EVENT_COMMENTS_VIEW_OWN, permissions.EVENT_COMMENTS_VIEW_ASSIGNED), eventCommentHandler.GetComment)
			protected.PUT("/event-comments/:id", middleware.RequirePermission(permissions.EVENT_COMMENTS_UPDATE_AUTHORED), eventCommentHandler.UpdateComment)
			protected.DELETE("/event-comments/:id", middleware.RequirePermission(permissions.EVENT_COMMENTS_DELETE_AUTHORED), eventCommentHandler.DeleteComment)
		}
	}

	return router
}
