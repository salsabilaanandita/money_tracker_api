package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingsGoal struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	Name          string        `gorm:"not null" json:"name"`
	TargetAmount  float64       `gorm:"not null" json:"target_amount"`
	CurrentAmount float64       `gorm:"default:0" json:"current_amount"`
	TargetDate    *time.Time    `json:"target_date"`
	CreatedAt     time.Time     `json:"created_at"`
	SavingEntries []SavingEntry `gorm:"foreignKey:SavingGoalID;constraint:OnDelete:CASCADE" json:"saving_entries,omitempty"`
}

func (sg *SavingsGoal) BeforeCreate(tx *gorm.DB) (err error) {
	if sg.ID == uuid.Nil {
		sg.ID = uuid.New()
	}
	return
}
