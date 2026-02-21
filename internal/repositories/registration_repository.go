package repositories

import (
	"context"

	"event_registration/internal/models"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	Create(ctx context.Context, registration *models.Registration) error
	FindByEventAndUser(ctx context.Context, eventID, userID string) (*models.Registration, error)
	FindByEvent(ctx context.Context, eventID string) ([]models.Registration, error)
	FindByUser(ctx context.Context, userID string) ([]models.Registration, error)
	UpdateStatus(ctx context.Context, registrationID string, status models.RegistrationStatus) error
	CountByEventAndStatus(ctx context.Context, eventID string, status models.RegistrationStatus) (int64, error)
	// For transaction purposes, we might need a way to pass the DB instance
	WithTx(tx *gorm.DB) RegistrationRepository
}

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) WithTx(tx *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: tx}
}

func (r *registrationRepository) Create(ctx context.Context, registration *models.Registration) error {
	return r.db.WithContext(ctx).Create(registration).Error
}

func (r *registrationRepository) FindByEventAndUser(ctx context.Context, eventID, userID string) (*models.Registration, error) {
	var registration models.Registration
	err := r.db.WithContext(ctx).Where("event_id = ? AND user_id = ?", eventID, userID).First(&registration).Error
	if err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *registrationRepository) FindByEvent(ctx context.Context, eventID string) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.WithContext(ctx).Preload("User").Where("event_id = ?", eventID).Find(&registrations).Error
	return registrations, err
}

func (r *registrationRepository) FindByUser(ctx context.Context, userID string) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.WithContext(ctx).Preload("Event").Where("user_id = ?", userID).Find(&registrations).Error
	return registrations, err
}

func (r *registrationRepository) UpdateStatus(ctx context.Context, registrationID string, status models.RegistrationStatus) error {
	return r.db.WithContext(ctx).Model(&models.Registration{}).Where("id = ?", registrationID).Update("status", status).Error
}

func (r *registrationRepository) CountByEventAndStatus(ctx context.Context, eventID string, status models.RegistrationStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Registration{}).Where("event_id = ? AND status = ?", eventID, status).Count(&count).Error
	return count, err
}
