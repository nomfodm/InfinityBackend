package entity

import "time"

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Token     string
	ExpiresAt time.Time
}
