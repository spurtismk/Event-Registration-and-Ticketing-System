package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "DRAFT"
	EventStatusPublished EventStatus = "PUBLISHED"
	EventStatusCancelled EventStatus = "CANCELLED"
)

type Event struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title          string      `gorm:"not null" json:"title"`
	Description    string      `json:"description"`
	Location       string      `json:"location"`
	EventDate      time.Time   `gorm:"not null" json:"event_date"`
	Capacity       int         `gorm:"not null;check:capacity >= 0" json:"capacity"`
	SeatsRemaining int         `gorm:"not null;check:seats_remaining >= 0" json:"seats_remaining"`
	OrganizerID    uuid.UUID   `gorm:"type:uuid;not null" json:"organizer_id"`
	Status         EventStatus `gorm:"type:varchar(20);not null;default:'DRAFT'" json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	Organizer     User           `gorm:"foreignKey:OrganizerID;references:ID" json:"organizer,omitempty"`
	Registrations []Registration `gorm:"foreignKey:EventID" json:"registrations,omitempty"`
	WaitlistItems []Waitlist     `gorm:"foreignKey:EventID" json:"waitlist_items,omitempty"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	// Verify constraints initially
	if e.SeatsRemaining > e.Capacity {
		e.SeatsRemaining = e.Capacity
	}
	return
}
