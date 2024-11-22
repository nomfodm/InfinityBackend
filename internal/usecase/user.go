package usecase

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
	"github.com/nomfodm/InfinityBackend/internal/utils"
	"io"
	"mime/multipart"
	"strings"
)

type UserUseCaseImpl struct {
	textureRepo repository.TextureRepository
	userRepo    repository.UserRepository
}

func NewUserUseCaseImpl(textureRepo repository.TextureRepository, userRepo repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{textureRepo: textureRepo, userRepo: userRepo}
}

func (uc *UserUseCaseImpl) SkinCapeHashes(user entity.User) (*string, *string) {
	var skinHash, capeHash *string

	skinID := user.SkinID
	capeID := user.CapeID

	if skinID != nil {
		skin, err := uc.textureRepo.SkinByID(*skinID)
		if err != nil {
			skinHash = nil
		} else {
			skinHash = &skin.TextureHash
		}
	} else {
		skinHash = nil
	}

	if capeID != nil {
		cape, err := uc.textureRepo.CapeByID(*capeID)
		if err != nil {
			capeHash = nil
		} else {
			capeHash = &cape.TextureHash
		}
	} else {
		capeHash = nil
	}

	return skinHash, capeHash

}

func (uc *UserUseCaseImpl) UploadSkin(user entity.User, skinFileHeader multipart.FileHeader) (string, error) {
	filename := skinFileHeader.Filename
	if !strings.HasSuffix(filename, ".png") {
		return "", errors.New("texture file type must be png")
	}

	skinFile, err := skinFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer skinFile.Close()

	skinFileBuffer := make([]byte, skinFileHeader.Size)
	skinFile.Read(skinFileBuffer)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, bytes.NewReader(skinFileBuffer)); err != nil {
		return "", err
	}
	fileHash := hasher.Sum(nil)
	fileHashString := fmt.Sprintf("%x", fileHash)

	skinInDB, err := uc.textureRepo.SkinByHash(fileHashString)
	if err == nil {
		err := uc.textureRepo.SetSkinToUser(user.ID, skinInDB.ID)
		return fileHashString, err
	}

	err = utils.UploadImagePNGToS3(fmt.Sprintf("textures/%s", fileHashString), skinFileBuffer)
	if err != nil {
		return "", err
	}

	avatarBuffer, err := utils.CropAvatarOutOfSkin(skinFileHeader)
	if err != nil {
		return "", err
	}

	err = utils.UploadImagePNGToS3(fmt.Sprintf("textures/avatars/%s", fileHashString), avatarBuffer)
	if err != nil {
		return "", err
	}

	skinID, err := uc.textureRepo.CreateSkin(fileHashString)
	if err != nil {
		return "", err
	}

	err = uc.textureRepo.SetSkinToUser(user.ID, skinID)
	return fileHashString, err
}

func (uc *UserUseCaseImpl) UploadCape(user entity.User, capeFileHeader multipart.FileHeader) (string, error) {
	filename := capeFileHeader.Filename
	if !strings.HasSuffix(filename, ".png") {
		return "", errors.New("texture file type must be png")
	}

	capeFile, err := capeFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer capeFile.Close()

	fileBuffer := make([]byte, capeFileHeader.Size)
	capeFile.Read(fileBuffer)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, bytes.NewReader(fileBuffer)); err != nil {
		return "", err
	}
	fileHash := hasher.Sum(nil)
	fileHashString := fmt.Sprintf("%x", fileHash)

	capeInDB, err := uc.textureRepo.CapeByHash(fileHashString)
	if err == nil {
		err := uc.textureRepo.SetCapeToUser(user.ID, capeInDB.ID)
		return fileHashString, err
	}

	err = utils.UploadImagePNGToS3(fmt.Sprintf("textures/%x", fileHash), fileBuffer)
	if err != nil {
		return "", err
	}

	capeID, err := uc.textureRepo.CreateCape(fileHashString)
	if err != nil {
		return "", err
	}

	err = uc.textureRepo.SetCapeToUser(user.ID, capeID)
	return fileHashString, err
}

func (uc *UserUseCaseImpl) ChangeNickname(user entity.User, newUsername string) error {
	return uc.userRepo.ChangeNickname(user.ID, newUsername)
}
