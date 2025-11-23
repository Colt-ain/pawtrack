package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/middleware"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/service"
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

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
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
			protected.POST("/events", eventHandler.CreateEvent)
			protected.GET("/events", eventHandler.ListEvents)
			protected.GET("/events/:id", eventHandler.GetEvent)
			protected.DELETE("/events/:id", eventHandler.DeleteEvent)

			// Dogs - create requires owner role, others just authentication
			protected.POST("/dogs", middleware.RequireRole(models.RoleOwner), dogHandler.CreateDog)
			protected.GET("/dogs", dogHandler.ListDogs)
			protected.GET("/dogs/:id", dogHandler.GetDog)
			protected.PUT("/dogs/:id", dogHandler.UpdateDog)
			protected.DELETE("/dogs/:id", dogHandler.DeleteDog)

			// Users - require authentication
			protected.GET("/users", userHandler.ListUsers)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)

			// Consultants - require authentication
			protected.PUT("/consultants/profile", consultantHandler.UpdateProfile)
			protected.GET("/consultants", consultantHandler.SearchConsultants)
			protected.GET("/consultants/:id", consultantHandler.GetProfile)
			protected.POST("/consultants/:id/invite", consultantHandler.InviteConsultant)

			// Invites - require authentication
			protected.POST("/invites/accept", consultantHandler.AcceptInvite)

			// Consultant Notes - require authentication
			protected.POST("/consultant-notes", consultantNoteHandler.CreateNote)
			protected.GET("/consultant-notes", consultantNoteHandler.ListNotes)
			protected.GET("/consultant-notes/:id", consultantNoteHandler.GetNote)
			protected.PUT("/consultant-notes/:id", consultantNoteHandler.UpdateNote)
			protected.DELETE("/consultant-notes/:id", consultantNoteHandler.DeleteNote)

			// Event Comments - require authentication
			protected.POST("/events/:id/comments", eventCommentHandler.CreateComment)
			protected.GET("/events/:id/comments", eventCommentHandler.ListComments)
			protected.GET("/event-comments/:id", eventCommentHandler.GetComment)
			protected.PUT("/event-comments/:id", eventCommentHandler.UpdateComment)
			protected.DELETE("/event-comments/:id", eventCommentHandler.DeleteComment)
		}
	}

	return router
}
