package repositories

import (
	"context"

	"event_registration/internal/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.Event) error
	FindByID(ctx context.Context, id string) (*models.Event, error)
	FindAll(ctx context.Context, status models.EventStatus) ([]models.Event, error)
	FindByOrganizer(ctx context.Context, organizerID string) ([]models.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*models.Event, error) {
	var event models.Event
	err := r.db.WithContext(ctx).Preload("Organizer").Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindAll(ctx context.Context, status models.EventStatus) ([]models.Event, error) {
	var events []models.Event
	query := r.db.WithContext(ctx).Preload("Organizer")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&events).Error
	return events, err
}

func (r *eventRepository) FindByOrganizer(ctx context.Context, organizerID string) ([]models.Event, error) {
	var events []models.Event
	err := r.db.WithContext(ctx).Preload("Organizer").Where("organizer_id = ?", organizerID).Find(&events).Error
	return events, err
}
