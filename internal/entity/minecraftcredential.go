package entity

import "github.com/google/uuid"

type MinecraftCredential struct {
	ID          uint `gorm:"primaryKey"`
	AccessToken string
	Username    string
	UUID        uuid.UUID
	ServerID    string
}
