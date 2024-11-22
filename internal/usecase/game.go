package usecase

import (
	"errors"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
)

type GameUseCaseImpl struct {
	gameRepo    repository.GameRepository
	textureRepo repository.TextureRepository
}

func NewGameUseCaseImpl(gameRepo repository.GameRepository,
	textureRepo repository.TextureRepository) *GameUseCaseImpl {
	return &GameUseCaseImpl{gameRepo: gameRepo, textureRepo: textureRepo}
}

var (
	ErrUserWithProvidedDetailsNotFound = errors.New("user with provided details not found")
)

func (uc *GameUseCaseImpl) Launcher(user entity.User) (username, uuid, accessToken string, err error) {
	userID := user.ID
	generatedAccessToken, err := uc.gameRepo.GenerateAccessTokenForUserAndSave(userID)
	return user.MinecraftCredential.Username, user.MinecraftCredential.UUID.String(), generatedAccessToken, err
}

func (uc *GameUseCaseImpl) Join(accessToken, uuid, serverID string) error {
	user, err := uc.gameRepo.UserByAccessTokenAndUUID(accessToken, uuid)
	if err != nil {
		return ErrUserWithProvidedDetailsNotFound
	}

	err = uc.gameRepo.ApplyServerIDToUser(user.ID, serverID)
	return err
}

func (uc *GameUseCaseImpl) HasJoined(username, serverID string) (entity.User, error) {
	user, err := uc.gameRepo.UserByUsernameAndServerID(username, serverID)
	return user, err
}

func (uc *GameUseCaseImpl) Profile(uuid string) (entity.User, *entity.Skin, *entity.Cape, error) {
	var skin *entity.Skin
	var cape *entity.Cape

	user, err := uc.gameRepo.UserByUUID(uuid)
	if err != nil {
		return user, skin, cape, err
	}

	if user.SkinID != nil {
		skinWithProvidedID, err := uc.textureRepo.SkinByID(*user.SkinID)
		skin = &skinWithProvidedID
		if err != nil {
			skin = nil
		}
	}

	if user.CapeID != nil {
		capeWithProvidedID, err := uc.textureRepo.CapeByID(*user.CapeID)
		cape = &capeWithProvidedID
		if err != nil {
			cape = nil
		}
	}

	return user, skin, cape, nil
}
