package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ActorID    uuid.UUID `gorm:"type:uuid;not null" json:"actor_id"`
	Action     string    `gorm:"not null" json:"action"`
	EntityType string    `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	Timestamp  time.Time `json:"timestamp"`

	Actor User `gorm:"foreignKey:ActorID;references:ID" json:"actor,omitempty"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
