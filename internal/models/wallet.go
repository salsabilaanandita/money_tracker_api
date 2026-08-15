package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`               // misal: "Cash", "BCA"
	Type      string    `gorm:"not null" json:"type"`               // bank, cash, e-wallet
	Balance   float64   `gorm:"default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}
