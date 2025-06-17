package usecase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
	"github.com/nomfodm/InfinityBackend/internal/utils"
	"io"
	"mime/multipart"
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
	if skinFileHeader.Header.Get("content-type") != "image/png" {
		return "", errors.New("skin file must be png")
	}

	skinFile, err := skinFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer skinFile.Close()

	var skinFileBuffer bytes.Buffer
	tee := io.TeeReader(skinFile, &skinFileBuffer)

	if err := utils.ValidateSkin(tee); err != nil {
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, tee); err != nil {
		return "", err
	}
	fileHash := hasher.Sum(nil)
	fileHashString := hex.EncodeToString(fileHash)

	skinInDB, err := uc.textureRepo.SkinByHash(fileHashString)
	if err == nil {
		err := uc.textureRepo.SetSkinToUser(user.ID, skinInDB.ID)
		return fileHashString, err
	}

	headBuffer, err := utils.RenderHeadOutOfSkin(skinFileHeader)
	if err != nil {
		return "", err
	}

	err = utils.UploadImagePNGToS3(fmt.Sprintf("textures/%s", fileHashString), skinFileBuffer.Bytes())
	if err != nil {
		return "", err
	}

	err = utils.UploadImagePNGToS3(fmt.Sprintf("textures/avatars/%s", fileHashString), headBuffer)
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
	if capeFileHeader.Header.Get("content-type") != "image/png" {
		return "", errors.New("skin file must be png")
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
