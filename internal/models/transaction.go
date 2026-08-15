package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WalletID    uuid.UUID `gorm:"type:uuid;not null" json:"wallet_id"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Type        string    `gorm:"not null" json:"type"` // income / expense
	Description string    `json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
