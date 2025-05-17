package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	Expires        time.Time
	SubscriptionID uint
	CreatedAt      time.Time
}
