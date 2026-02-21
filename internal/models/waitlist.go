package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Waitlist struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_waitlist_user_event" json:"user_id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_waitlist_user_event" json:"event_id"`
	Position  int       `gorm:"not null" json:"position"`
	CreatedAt time.Time `json:"created_at"`

	User  User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Event Event `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
}

func (w *Waitlist) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}
