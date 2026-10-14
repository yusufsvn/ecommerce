package models

import (
	"time"

	"github.com/google/uuid"
)

type Refreshtoken struct {
	Token_id   uuid.UUID
	User_id    uuid.UUID
	Token_hash string
	family_id  uuid.UUID
	Is_used    bool
	Is_revoked bool
	Expires_at time.Time
	Created_At time.Time
	Updated_At time.Time
}
