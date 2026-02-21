package db

import (
	"log"
	"time"

	"event_registration/internal/models"
	"event_registration/internal/utils"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already contains data. Skipping seed.")
		return
	}

	log.Println("Seeding Database with sample Users and Events...")

	// Create Users
	adminPwd, _ := utils.HashPassword("admin123")
	orgPwd, _ := utils.HashPassword("org123")
	userPwd, _ := utils.HashPassword("user123")

	admin := models.User{Name: "Super Admin", Email: "admin@eventbrite.local", PasswordHash: adminPwd, Role: models.RoleAdmin}
	org1 := models.User{Name: "TechCorp Conferences", Email: "events@techcorp.com", PasswordHash: orgPwd, Role: models.RoleOrganizer}
	user1 := models.User{Name: "Alice Student", Email: "alice@student.com", PasswordHash: userPwd, Role: models.RoleAudience}
	user2 := models.User{Name: "Bob Engineer", Email: "bob@engineer.com", PasswordHash: userPwd, Role: models.RoleAudience}

	db.Create(&admin)
	db.Create(&org1)
	db.Create(&user1)
	db.Create(&user2)

	// Create Events
	e1 := models.Event{
		Title:          "Golang Microservices Summit 2026",
		Description:    "Advanced patterns for building highly scalable systems using Go and Docker.",
		Location:       "Convention Center A",
		EventDate:      time.Now().Add(24 * time.Hour),
		Capacity:       10, // Very small capacity to test concurrency easily
		SeatsRemaining: 10,
		OrganizerID:    org1.ID,
		Status:         models.EventStatusPublished,
	}

	e2 := models.Event{
		Title:          "Introduction to PostgreSQL Locking",
		Description:    "Learn how SELECT FOR UPDATE prevents race conditions inside explicit database transactions.",
		Location:       "Online Webinar",
		EventDate:      time.Now().Add(72 * time.Hour),
		Capacity:       5, // Small capacity
		SeatsRemaining: 5,
		OrganizerID:    org1.ID,
		Status:         models.EventStatusPublished,
	}

	e3 := models.Event{
		Title:          "Draft Event (Not Visible)",
		Description:    "This event is still being planned.",
		Location:       "TBD",
		EventDate:      time.Now().Add(100 * time.Hour),
		Capacity:       100,
		SeatsRemaining: 100,
		OrganizerID:    org1.ID,
		Status:         models.EventStatusDraft,
	}

	db.Create(&e1)
	db.Create(&e2)
	db.Create(&e3)

	log.Println("Database Seeding Completed Successfully!")
}
