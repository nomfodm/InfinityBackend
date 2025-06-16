package repository

import "github.com/nomfodm/InfinityBackend/internal/entity"

type UserRepository interface {
	Create(username string, email string, passwordHash string) error
	ByUsername(username string) (entity.User, error)
	ByEmail(email string) (entity.User, error)
	ByID(userID uint) (entity.User, error)
	FindRefreshToken(refreshToken string) (entity.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
	SaveRefreshToken(userID uint, refreshToken string) error
	Activate(userID uint) error
	ChangeNickname(userID uint, newUsername string) error
}

type TextureRepository interface {
	SkinByID(skinID uint) (entity.Skin, error)
	CapeByID(capeID uint) (entity.Cape, error)
	SkinByHash(hash string) (entity.Skin, error)
	CapeByHash(hash string) (entity.Cape, error)
	CreateSkin(hash string) (uint, error)
	CreateCape(hash string) (uint, error)
	SetSkinToUser(userID, skinID uint) error
	SetCapeToUser(userID, capeID uint) error
}

type GameRepository interface {
	GenerateAccessTokenForUserAndSave(userID uint) (string, error)
	UserByAccessTokenAndUUID(accessToken, uuid string) (entity.User, error)
	UserByUsernameAndServerID(username, serverID string) (entity.User, error)
	UserByUUID(uuid string) (entity.User, error)
	ApplyServerIDToUser(userID uint, serverID string) error
}

type LauncherRepository interface {
	LatestLauncherVersion() (entity.LauncherVersion, error)
	CreateNewLauncherVersion(version entity.LauncherVersion) error
	LastMandatoryVersion() (entity.LauncherVersion, error)
}

type ServerStatusRepository interface {
	InitServerStatus() error
	CurrentServerStatus() (entity.ServerStatus, error)
	SetServerStatus(newStatus int) error
}
