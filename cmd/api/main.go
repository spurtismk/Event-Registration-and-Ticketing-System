package main

import (
	"log"

	"event_registration/internal/config"
	"event_registration/internal/db"
	"event_registration/internal/handlers"
	"event_registration/internal/repositories"
	"event_registration/internal/router"
	"event_registration/internal/services"
)

func main() {
	log.Println("Loading Configuration...")
	cfg := config.LoadConfig()

	log.Println("Initializing Database...")
	database := db.InitDB(cfg)

	// Seed Sample Data Automatically
	db.SeedDatabase(database)

	// Repositories
	userRepo := repositories.NewUserRepository(database)
	eventRepo := repositories.NewEventRepository(database)
	regRepo := repositories.NewRegistrationRepository(database)
	waitRepo := repositories.NewWaitlistRepository(database)

	// Services
	authService := services.NewAuthService(userRepo)
	eventService := services.NewEventService(eventRepo)
	regService := services.NewRegistrationService(database, regRepo, waitRepo, eventRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, cfg.JWTSecret)
	eventHandler := handlers.NewEventHandler(eventService, regService)
	organizerHandler := handlers.NewOrganizerHandler(eventService, regService)
	adminHandler := handlers.NewAdminHandler(database, regService, eventService)

	// Router
	log.Println("Setting up Router...")
	r := router.SetupRouter(
		cfg.JWTSecret,
		authHandler,
		eventHandler,
		organizerHandler,
		adminHandler,
	)

	log.Printf("Starting Server on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
