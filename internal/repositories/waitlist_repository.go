package repositories

import (
	"context"

	"event_registration/internal/models"
	"gorm.io/gorm"
)

type WaitlistRepository interface {
	Create(ctx context.Context, waitlist *models.Waitlist) error
	GetNextInLine(ctx context.Context, eventID string) (*models.Waitlist, error)
	Delete(ctx context.Context, waitlistID string) error
	CountByEvent(ctx context.Context, eventID string) (int64, error)
	FindByEventAndUser(ctx context.Context, eventID, userID string) (*models.Waitlist, error)
	WithTx(tx *gorm.DB) WaitlistRepository
}

type waitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository(db *gorm.DB) WaitlistRepository {
	return &waitlistRepository{db: db}
}

func (r *waitlistRepository) WithTx(tx *gorm.DB) WaitlistRepository {
	return &waitlistRepository{db: tx}
}

func (r *waitlistRepository) Create(ctx context.Context, waitlist *models.Waitlist) error {
	return r.db.WithContext(ctx).Create(waitlist).Error
}

func (r *waitlistRepository) GetNextInLine(ctx context.Context, eventID string) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Order("position asc").First(&waitlist).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}

func (r *waitlistRepository) Delete(ctx context.Context, waitlistID string) error {
	return r.db.WithContext(ctx).Where("id = ?", waitlistID).Delete(&models.Waitlist{}).Error
}

func (r *waitlistRepository) CountByEvent(ctx context.Context, eventID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Waitlist{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}

func (r *waitlistRepository) FindByEventAndUser(ctx context.Context, eventID, userID string) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.WithContext(ctx).Where("event_id = ? AND user_id = ?", eventID, userID).First(&waitlist).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}
