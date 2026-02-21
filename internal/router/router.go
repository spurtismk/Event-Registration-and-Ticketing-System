package router

import (
	"event_registration/internal/handlers"
	"event_registration/internal/middleware"
	"event_registration/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	jwtSecret string,
	authHandler *handlers.AuthHandler,
	eventHandler *handlers.EventHandler,
	organizerHandler *handlers.OrganizerHandler,
	adminHandler *handlers.AdminHandler,
) *gin.Engine {
	r := gin.Default()

	// Serve Static Frontend
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Auth Routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Audience / General Event Browsing
	events := r.Group("/events")
	events.Use(middleware.AuthRequired(jwtSecret))
	{
		events.GET("", eventHandler.ListEvents)
		events.GET("/:id", eventHandler.GetEvent)
		events.POST("/:id/register", eventHandler.RegisterForEvent)
		events.POST("/registrations/:registration_id/cancel", eventHandler.CancelRegistration)
	}

	// Organizer Routes
	organizer := r.Group("/organizer")
	organizer.Use(middleware.AuthRequired(jwtSecret), middleware.RoleRequired(models.RoleOrganizer, models.RoleAdmin))
	{
		organizer.POST("/events", organizerHandler.CreateEvent)
		organizer.GET("/events", organizerHandler.ListMyEvents)
		organizer.POST("/events/:id/publish", organizerHandler.PublishEvent)
		organizer.POST("/events/:id/cancel", organizerHandler.CancelEvent)
		organizer.GET("/events/:id/analytics", organizerHandler.GetAnalytics)
	}

	// Admin Routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(jwtSecret), middleware.RoleRequired(models.RoleAdmin))
	{
		admin.POST("/events/:id/simulate", adminHandler.SimulateConcurrency)
	}

	return r
}
