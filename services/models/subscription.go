package models

import "time"

type Subscription struct {
	ID        uint
	Email     string `gorm:"unique"`
	City      string
	Frequency string
	Confirmed bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Tokens    []Token
}
