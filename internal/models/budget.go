package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	AmountLimit float64   `gorm:"not null" json:"amount_limit"`
	Month       int       `gorm:"not null" json:"month"` // 1-12
	Year        int       `gorm:"not null" json:"year"`
	CreatedAt   time.Time `json:"created_at"`
}

func (b *Budget) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
