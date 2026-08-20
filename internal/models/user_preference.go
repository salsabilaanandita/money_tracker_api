package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserPreference menyimpan pengaturan tampilan, notifikasi, dan privasi milik user.
// Satu user hanya punya satu baris preferensi (relasi 1-1 dengan User).
// Nama field JSON di sini HARUS sama dengan yang dikirim frontend (lihat type
// Preferences di halaman Preferensi Next.js).
type UserPreference struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`

	// --- Tampilan ---
	Theme string `gorm:"not null;default:'system'" json:"theme"` // light, dark, system

	// --- Notifikasi ---
	BudgetNotification      bool `gorm:"not null;default:true" json:"budget_notification"`
	TransactionNotification bool `gorm:"not null;default:true" json:"transaction_notification"`
	SavingsNotification     bool `gorm:"not null;default:true" json:"savings_notification"`
	ReminderNotification    bool `gorm:"not null;default:false" json:"reminder_notification"`

	// --- Privasi ---
	PrivateMode bool `gorm:"not null;default:false" json:"private_mode"`
	HideBalance bool `gorm:"not null;default:false" json:"hide_balance"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *UserPreference) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}