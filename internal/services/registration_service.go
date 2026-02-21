package services

import (
	"context"
	"errors"
	"time"

	"event_registration/internal/models"
	"event_registration/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RegistrationService interface {
	BookEvent(ctx context.Context, userID, eventID string) (*models.Registration, *models.Waitlist, error)
	CancelRegistration(ctx context.Context, userID, registrationID string) error
	GetOrganizerAnalytics(ctx context.Context, organizerID, eventID string) (map[string]interface{}, error)
}

type registrationService struct {
	db        *gorm.DB
	regRepo   repositories.RegistrationRepository
	waitRepo  repositories.WaitlistRepository
	eventRepo repositories.EventRepository
}

func NewRegistrationService(db *gorm.DB, regRepo repositories.RegistrationRepository, waitRepo repositories.WaitlistRepository, eventRepo repositories.EventRepository) RegistrationService {
	return &registrationService{
		db:        db,
		regRepo:   regRepo,
		waitRepo:  waitRepo,
		eventRepo: eventRepo,
	}
}

// BookEvent contains the core concurrency-safe logic
func (s *registrationService) BookEvent(ctx context.Context, userID, eventID string) (*models.Registration, *models.Waitlist, error) {
	var finalReg *models.Registration
	var finalWaitlist *models.Waitlist

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock the event row using SELECT ... FOR UPDATE
		var event models.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", eventID).First(&event).Error; err != nil {
			return errors.New("event not found")
		}

		// 2. Validate Event criteria
		if event.Status != models.EventStatusPublished {
			return errors.New("event is not published")
		}
		if event.EventDate.Before(time.Now()) {
			return errors.New("event has already passed")
		}

		// 3. User Parsing
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		// Check if user already registered or waitlisted
		var existingReg models.Registration
		if err := tx.Where("event_id = ? AND user_id = ?", eventID, userID).First(&existingReg).Error; err == nil {
			if existingReg.Status == models.RegistrationStatusConfirmed {
				return errors.New("already registered for this event")
			}
		}

		var existingWaitlist models.Waitlist
		if err := tx.Where("event_id = ? AND user_id = ?", eventID, userID).First(&existingWaitlist).Error; err == nil {
			return errors.New("already on waitlist for this event")
		}

		// 4. Concurrency Logic / Distribution
		if event.SeatsRemaining > 0 {
			// Seat available
			newReg := &models.Registration{
				UserID:  userUUID,
				EventID: event.ID,
				Status:  models.RegistrationStatusConfirmed,
			}
			if err := tx.Create(newReg).Error; err != nil {
				return err
			}

			// Decrement seats
			event.SeatsRemaining -= 1
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
			finalReg = newReg
		} else {
			// Waitlist
			var count int64
			tx.Model(&models.Waitlist{}).Where("event_id = ?", event.ID).Count(&count)

			newWaitlist := &models.Waitlist{
				UserID:   userUUID,
				EventID:  event.ID,
				Position: int(count) + 1,
			}
			if err := tx.Create(newWaitlist).Error; err != nil {
				return err
			}
			finalWaitlist = newWaitlist
		}

		return nil // Commit transaction
	})

	return finalReg, finalWaitlist, err
}

// CancelRegistration cancels the reg, increments seat or polls from waitlist
func (s *registrationService) CancelRegistration(ctx context.Context, userID, registrationID string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reg models.Registration
		if err := tx.Where("id = ?", registrationID).First(&reg).Error; err != nil {
			return errors.New("registration not found")
		}

		if reg.UserID.String() != userID {
			return errors.New("unauthorized to cancel this registration")
		}

		if reg.Status == models.RegistrationStatusCancelled {
			return errors.New("already cancelled")
		}

		// Lock event
		var event models.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", reg.EventID).First(&event).Error; err != nil {
			return err
		}

		// Cancel registration
		reg.Status = models.RegistrationStatusCancelled
		if err := tx.Save(&reg).Error; err != nil {
			return err
		}

		// Check waitlist
		var nextUser models.Waitlist
		err := tx.Where("event_id = ?", event.ID).Order("position asc").First(&nextUser).Error
		if err == nil { // Waitlist user found
			// Promote user to Registration
			newReg := &models.Registration{
				UserID:  nextUser.UserID,
				EventID: event.ID,
				Status:  models.RegistrationStatusConfirmed,
			}
			if err := tx.Create(newReg).Error; err != nil {
				return err
			}

			// Remove from waitlist
			if err := tx.Delete(&nextUser).Error; err != nil {
				return err
			}
			// Seats remaining does not change because it's transferred
		} else {
			// No waitlist, increment seats
			event.SeatsRemaining += 1
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *registrationService) GetOrganizerAnalytics(ctx context.Context, organizerID, eventID string) (map[string]interface{}, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.OrganizerID.String() != organizerID {
		return nil, errors.New("unauthorized access to analytics")
	}

	confirmedCount, _ := s.regRepo.CountByEventAndStatus(ctx, eventID, models.RegistrationStatusConfirmed)
	var waitlistCount int64
	s.db.WithContext(ctx).Model(&models.Waitlist{}).Where("event_id = ?", eventID).Count(&waitlistCount)

	totalRegistrations := confirmedCount + waitlistCount
	var seatsFilledPercentage float64
	if event.Capacity > 0 {
		seatsFilledPercentage = float64(event.Capacity-event.SeatsRemaining) / float64(event.Capacity) * 100
	}

	return map[string]interface{}{
		"total_registrations":     totalRegistrations,
		"confirmed_count":         confirmedCount,
		"waitlist_count":          waitlistCount,
		"seats_filled_percentage": seatsFilledPercentage,
	}, nil
}
