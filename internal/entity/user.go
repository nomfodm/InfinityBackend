package entity

import "time"

type User struct {
	ID           uint `gorm:"primaryKey"`
	Username     string
	Email        string
	PasswordHash string
	Active       bool

	RefreshTokens []RefreshToken

	MinecraftCredentialID uint
	MinecraftCredential   MinecraftCredential

	SkinID *uint
	CapeID *uint

	RegisteredAt time.Time
}
