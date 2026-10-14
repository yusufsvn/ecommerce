package models

import (
	"time"

	"github.com/google/uuid"
)

// Address tablosu
type Address struct {
	AddressID   uuid.UUID  `db:"address_id" json:"address_id"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	Name        string     `db:"name" json:"name"` // Adres başlığı (Ev, İş vb.)
	PostCode    string     `db:"post_code" json:"post_code"`
	Country     string     `db:"country" json:"country"`
	City        string     `db:"city" json:"city"`
	District    string     `db:"district" json:"district"`
	FullAddress string     `db:"full_address" json:"full_address"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}
