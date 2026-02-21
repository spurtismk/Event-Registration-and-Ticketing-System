package services

import (
	"context"
	"errors"

	"event_registration/internal/models"
	"event_registration/internal/repositories"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, organizerID string, event *models.Event) error
	PublishEvent(ctx context.Context, organizerID, eventID string) error
	CancelEvent(ctx context.Context, organizerID, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*models.Event, error)
	ListPublishedEvents(ctx context.Context) ([]models.Event, error)
	ListOrganizerEvents(ctx context.Context, organizerID string) ([]models.Event, error)
}

type eventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) CreateEvent(ctx context.Context, organizerID string, event *models.Event) error {
	orgUUID, err := uuid.Parse(organizerID)
	if err != nil {
		return errors.New("invalid organizer ID")
	}
	event.OrganizerID = orgUUID
	event.SeatsRemaining = event.Capacity
	event.Status = models.EventStatusDraft
	return s.eventRepo.Create(ctx, event)
}

func (s *eventService) PublishEvent(ctx context.Context, organizerID, eventID string) error {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.OrganizerID.String() != organizerID {
		return errors.New("unauthorized to publish this event")
	}

	if event.Status != models.EventStatusDraft {
		return errors.New("event is not in draft status")
	}

	event.Status = models.EventStatusPublished
	return s.eventRepo.Update(ctx, event)
}

func (s *eventService) CancelEvent(ctx context.Context, organizerID, eventID string) error {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.OrganizerID.String() != organizerID {
		return errors.New("unauthorized to cancel this event")
	}

	event.Status = models.EventStatusCancelled
	return s.eventRepo.Update(ctx, event)
}

func (s *eventService) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return s.eventRepo.FindByID(ctx, eventID)
}

func (s *eventService) ListPublishedEvents(ctx context.Context) ([]models.Event, error) {
	return s.eventRepo.FindAll(ctx, models.EventStatusPublished)
}

func (s *eventService) ListOrganizerEvents(ctx context.Context, organizerID string) ([]models.Event, error) {
	return s.eventRepo.FindByOrganizer(ctx, organizerID)
}
