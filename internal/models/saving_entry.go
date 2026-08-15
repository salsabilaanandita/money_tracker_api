package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingEntry struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SavingGoalID uuid.UUID `gorm:"type:uuid;not null" json:"saving_goal_id"`
	Amount       float64   `gorm:"not null" json:"amount"`
	Date         time.Time `gorm:"not null" json:"date"`
	CreatedAt    time.Time `json:"created_at"`
}

func (se *SavingEntry) BeforeCreate(tx *gorm.DB) (err error) {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return
}
