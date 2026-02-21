package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"event_registration/internal/models"
	"event_registration/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db           *gorm.DB
	regService   services.RegistrationService
	eventService services.EventService
}

func NewAdminHandler(db *gorm.DB, regService services.RegistrationService, eventService services.EventService) *AdminHandler {
	return &AdminHandler{
		db:           db,
		regService:   regService,
		eventService: eventService,
	}
}

// SimulateConcurrency hits the BookEvent method across N goroutines
func (h *AdminHandler) SimulateConcurrency(c *gin.Context) {
	eventID := c.Param("id")
	usersParam := c.Query("users")

	numUsers, err := strconv.Atoi(usersParam)
	if err != nil || numUsers <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter: users"})
		return
	}

	// 1. Check if event is valid
	ctx := context.Background()
	event, err := h.eventService.GetEvent(ctx, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// 2. We must create N dummy users to satisfy foreign keys for the simulation.
	// Normally we wouldn't do this in a production handler, but it's part of the assignment.
	var userIDs []string
	
	// Create a unique simulation batch string to prevent email collisions across runs
	simBatch := uuid.New().String()[:8]
	
	for i := 0; i < numUsers; i++ {
		user := &models.User{
			Name:         fmt.Sprintf("Sim User %d", i),
			Email:        fmt.Sprintf("sim%d_%s_%s@sim.local", i, simBatch, event.ID.String()[:4]), // totally unique
			PasswordHash: "not_needed",
			Role:         models.RoleAudience,
		}
		if err := h.db.WithContext(ctx).Create(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dummy users: " + err.Error()})
			return
		}
		userIDs = append(userIDs, user.ID.String())
	}

	// 3. Spawning Goroutines
	var successCount int64
	var waitlistCount int64
	var failCount int64

	var wg sync.WaitGroup
	wg.Add(numUsers)

	for _, uID := range userIDs {
		go func(userID string) {
			defer wg.Done()
			
			// Each goroutine represents a concurrent request
			_, waitlist, err := h.regService.BookEvent(context.Background(), userID, eventID)
			
			if err != nil {
				atomic.AddInt64(&failCount, 1)
			} else if waitlist != nil {
				atomic.AddInt64(&waitlistCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(uID)
	}

	wg.Wait()

	// 4. Fetch final remaining seats
	finalEvent, _ := h.eventService.GetEvent(ctx, eventID)

	c.JSON(http.StatusOK, gin.H{
		"simulation_results": gin.H{
			"total_attempted": numUsers,
			"success_count":   successCount,
			"waitlisted_count": waitlistCount,
			"failed_count":    failCount,
			"final_seats_remaining": finalEvent.SeatsRemaining,
		},
	})
}
