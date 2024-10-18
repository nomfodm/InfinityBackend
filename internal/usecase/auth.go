package usecase

import (
	"backend/internal/entity"
	"backend/internal/infrastructure/repository"
	"backend/internal/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthUseCaseImpl struct {
	repo repository.UserRepository
}

func NewAuthUseCaseImpl(repo repository.UserRepository) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{repo: repo}
}

var (
	ErrUsernameAlreadyInUse = errors.New("provided username is already in use")
	ErrEmailAlreadyInUse    = errors.New("provided email is already in use")
	ErrCantHashPassword     = errors.New("cannot perform password hashing")
	ErrInvalidPassword      = errors.New("invalid password")

	ErrUserNotFound         = errors.New("user not found in database")
	ErrCantSaveRefreshToken = errors.New("cant save refresh token")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")

	ErrCantParseUUID = errors.New("cannot parse uuid")
	ErrInvalidUUID   = errors.New("invalid uuid")
)

func (uc *AuthUseCaseImpl) SignUp(username string, email string, password string) error {
	if _, err := uc.repo.ByUsername(username); err == nil {
		return ErrUsernameAlreadyInUse
	}

	if _, err := uc.repo.ByEmail(email); err == nil {
		return ErrEmailAlreadyInUse
	}

	passwordHash, err := utils.HashStringToBcrypt(password)
	if err != nil {
		return ErrCantHashPassword
	}

	err = uc.repo.Create(username, email, passwordHash)
	return err
}

func (uc *AuthUseCaseImpl) generateTokensForUser(userID uint) (string, string, error) {
	accessToken, err := utils.GenerateJWTForUser(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken := utils.GenerateRefreshToken()
	err = uc.repo.SaveRefreshToken(userID, refreshToken)
	if err != nil {
		return "", "", ErrCantSaveRefreshToken
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUseCaseImpl) SignIn(username string, password string) (string, string, error) {
	user, err := uc.repo.ByUsername(username)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	if !utils.VerifyStringHash(user.PasswordHash, password) {
		return "", "", ErrInvalidPassword
	}

	userID := user.ID

	return uc.generateTokensForUser(userID)
}

func (uc *AuthUseCaseImpl) Refresh(refreshToken string) (string, string, error) {
	refreshTokenRow, err := uc.repo.FindRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrRefreshTokenNotFound
	}
	uc.repo.DeleteRefreshToken(refreshToken)

	if time.Now().After(refreshTokenRow.ExpiresAt) {
		return "", "", ErrRefreshTokenExpired
	}

	userID := refreshTokenRow.UserID

	_, err = uc.repo.ByID(userID)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	return uc.generateTokensForUser(userID)
}

func (uc *AuthUseCaseImpl) Logout(refreshToken string) error {
	return uc.repo.DeleteRefreshToken(refreshToken)
}

func (uc *AuthUseCaseImpl) Activate(userID uint, minecraftUUID string) error {
	user, err := uc.repo.ByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	parsedUUID, err := uuid.Parse(minecraftUUID)
	if err != nil {
		return ErrCantParseUUID
	}

	if user.MinecraftCredential.UUID != parsedUUID {
		return ErrInvalidUUID
	}

	return uc.repo.Activate(user.ID)
}

func (uc *AuthUseCaseImpl) RetrieveUserWithAccessToken(accessToken string) (entity.User, error) {
	userID, err := utils.ParseUserJWT(accessToken)
	if err != nil {
		return entity.User{}, err
	}

	user, err := uc.repo.ByID(userID)
	if err != nil {
		return user, ErrUserNotFound
	}

	return user, nil
}
