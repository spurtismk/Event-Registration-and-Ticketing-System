package handlers

import (
	"net/http"

	"event_registration/internal/services"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService services.EventService
	regService   services.RegistrationService
}

func NewEventHandler(eventService services.EventService, regService services.RegistrationService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		regService:   regService,
	}
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	events, err := h.eventService.ListPublishedEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	id := c.Param("id")
	event, err := h.eventService.GetEvent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (h *EventHandler) RegisterForEvent(c *gin.Context) {
	eventID := c.Param("id")
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(string)

	reg, waitlist, err := h.regService.BookEvent(c.Request.Context(), userID, eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if waitlist != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Event is full. Added to waitlist.",
			"waitlist": waitlist,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Successfully registered for event",
		"registration": reg,
	})
}

func (h *EventHandler) CancelRegistration(c *gin.Context) {
	regID := c.Param("registration_id")
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(string)

	err := h.regService.CancelRegistration(c.Request.Context(), userID, regID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration cancelled successfully"})
}
