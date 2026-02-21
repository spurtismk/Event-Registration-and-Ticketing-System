package handlers

import (
	"net/http"

	"event_registration/internal/models"
	"event_registration/internal/services"
	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct {
	eventService services.EventService
	regService   services.RegistrationService
}

func NewOrganizerHandler(eventService services.EventService, regService services.RegistrationService) *OrganizerHandler {
	return &OrganizerHandler{
		eventService: eventService,
		regService:   regService,
	}
}

func (h *OrganizerHandler) CreateEvent(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	organizerID := userIDVal.(string)

	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), organizerID, &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created in Draft status",
		"event":   event,
	})
}

func (h *OrganizerHandler) PublishEvent(c *gin.Context) {
	eventID := c.Param("id")
	userIDVal, _ := c.Get("userID")
	organizerID := userIDVal.(string)

	if err := h.eventService.PublishEvent(c.Request.Context(), organizerID, eventID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event published successfully"})
}

func (h *OrganizerHandler) CancelEvent(c *gin.Context) {
	eventID := c.Param("id")
	userIDVal, _ := c.Get("userID")
	organizerID := userIDVal.(string)

	if err := h.eventService.CancelEvent(c.Request.Context(), organizerID, eventID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event cancelled successfully"})
}

func (h *OrganizerHandler) ListMyEvents(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	organizerID := userIDVal.(string)

	events, err := h.eventService.ListOrganizerEvents(c.Request.Context(), organizerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *OrganizerHandler) GetAnalytics(c *gin.Context) {
	eventID := c.Param("id")
	userIDVal, _ := c.Get("userID")
	organizerID := userIDVal.(string)

	analytics, err := h.regService.GetOrganizerAnalytics(c.Request.Context(), organizerID, eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}
