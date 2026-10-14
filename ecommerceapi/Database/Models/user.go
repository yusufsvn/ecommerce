package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID  `db:"user_id" json:"user_id"`
	Password      string     `db:"password" json:"-"`
	Name          string     `db:"name" json:"name"`
	Surname       string     `db:"surname" json:"surname"`
	Phone         string     `db:"phone" json:"phone"`
	Email         string     `db:"email" json:"email"`
	UserRole      string     `db:"user_role" json:"user_role"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	OAuthProvider *string    `db:"oauth_provider" json:"oauth_provider"`
	OAuthID       *string    `db:"oauth_id" json:"oauth_id"`
}

