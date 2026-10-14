package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductCategory struct {
	CategoryID     uuid.UUID  `db:"category_id" json:"category_id"`
	CategoryName   string     `db:"category_name" json:"category_name"`
	CategoryEnumNo int        `db:"category_enum_no" json:"category_enum_no"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
}
