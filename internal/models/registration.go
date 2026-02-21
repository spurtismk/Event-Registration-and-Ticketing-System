package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrationStatus string

const (
	RegistrationStatusConfirmed RegistrationStatus = "CONFIRMED"
	RegistrationStatusCancelled RegistrationStatus = "CANCELLED"
)

type Registration struct {
	ID        uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex:idx_user_event" json:"user_id"`
	EventID   uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex:idx_user_event" json:"event_id"`
	Status    RegistrationStatus `gorm:"type:varchar(20);not null;default:'CONFIRMED'" json:"status"`
	CreatedAt time.Time          `json:"created_at"`

	User  User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Event Event `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
}

func (r *Registration) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
