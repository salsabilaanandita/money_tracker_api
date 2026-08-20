package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserPreference menyimpan pengaturan tampilan, notifikasi, dan privasi milik user.
// Satu user hanya punya satu baris preferensi (relasi 1-1 dengan User).
type UserPreference struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`

	// --- Tampilan ---
	Theme          string `gorm:"not null;default:'light'" json:"theme"`         // light, dark, system
	Language       string `gorm:"not null;default:'id'" json:"language"`         // id, en
	CurrencyFormat string `gorm:"not null;default:'IDR'" json:"currency_format"` // IDR, USD, dst

	// --- Notifikasi ---
	NotifTransaction     bool `gorm:"not null;default:true" json:"notif_transaction"`      // notifikasi transaksi baru
	NotifBudgetAlert     bool `gorm:"not null;default:true" json:"notif_budget_alert"`     // notifikasi budget mendekati/lewat batas
	NotifSavingsReminder bool `gorm:"not null;default:true" json:"notif_savings_reminder"` // pengingat target tabungan
	NotifEmail           bool `gorm:"not null;default:false" json:"notif_email"`           // notifikasi via email

	// --- Privasi ---
	HideBalance        bool `gorm:"not null;default:false" json:"hide_balance"`         // sembunyikan saldo di halaman utama
	BiometricLock      bool `gorm:"not null;default:false" json:"biometric_lock"`       // kunci aplikasi dengan biometrik
	ShareDataAnalytics bool `gorm:"not null;default:false" json:"share_data_analytics"` // izinkan berbagi data untuk analitik

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *UserPreference) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
