package usecase

import (
	"backend/internal/entity"
	"mime/multipart"
)

type AuthUseCase interface {
	SignUp(username string, email string, password string) error
	SignIn(username string, password string) (string, string, error)
	Logout(refreshToken string) error
	Refresh(refreshToken string) (string, string, error)
	Activate(userID uint, minecraftUUID string) error

	RetrieveUserWithAccessToken(accessToken string) (entity.User, error)
}

type UserUseCase interface {
	SkinCapeHashes(user entity.User) (*string, *string)
	UploadSkin(user entity.User, skinFileHeader multipart.FileHeader) (string, error)
	UploadCape(user entity.User, capeFileHeader multipart.FileHeader) (string, error)
	ChangeNickname(user entity.User, newUsername string) error
}

type GameUseCase interface {
	Launcher(user entity.User) (username, uuid, accessToken string, err error)
	Join(accessToken, uuid, serverID string) error
	HasJoined(username, serverID string) (entity.User, error)
	Profile(uuid string) (entity.User, *entity.Skin, *entity.Cape, error)
}
