package user

import (
	"time"
)

type meResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Active       bool      `json:"active"`
	RegisteredAt time.Time `json:"registeredAt"`

	Textures            textures            `json:"textures"`
	MinecraftCredential minecraftCredential `json:"minecraftCredential"`
}

type textures struct {
	SkinHash *string `json:"skinHash"`
	CapeHash *string `json:"capeHash"`
}

type minecraftCredential struct {
	Username string `json:"username"`
}
