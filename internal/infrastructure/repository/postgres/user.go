package postgres

import (
	"errors"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) Create(username string, email string, passwordHash string) error {
	userToCreate := entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		MinecraftCredential: entity.MinecraftCredential{
			Username: username,
			UUID:     uuid.New(),
		},
		RegisteredAt: time.Now(),
	}

	result := repo.db.Create(&userToCreate)
	return result.Error
}

func (repo *PostgresUserRepository) ByUsername(username string) (entity.User, error) {
	var user entity.User
	result := repo.db.Preload(clause.Associations).Where("username = ?", username).First(&user)
	return user, result.Error
}

func (repo *PostgresUserRepository) ByEmail(email string) (entity.User, error) {
	var user entity.User
	result := repo.db.Preload(clause.Associations).Where("email = ?", email).First(&user)
	return user, result.Error
}

func (repo *PostgresUserRepository) ByID(userID uint) (entity.User, error) {
	var user entity.User
	result := repo.db.Preload(clause.Associations).Where("id = ?", userID).First(&user)
	return user, result.Error
}

func (repo *PostgresUserRepository) FindRefreshToken(refreshToken string) (entity.RefreshToken, error) {
	var token entity.RefreshToken
	result := repo.db.Preload(clause.Associations).Where("token = ?", refreshToken).First(&token)
	return token, result.Error
}

func (repo *PostgresUserRepository) DeleteRefreshToken(refreshToken string) error {
	refreshTokenRow, err := repo.FindRefreshToken(refreshToken)
	if err != nil {
		return errors.New("refresh token not found")
	}
	result := repo.db.Delete(&refreshTokenRow)
	return result.Error
}

func (repo *PostgresUserRepository) SaveRefreshToken(userID uint, refreshToken string) error {
	refreshTokenLifetime, _ := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_LIFETIME"), 10, 32)
	refreshTokenToCreate := entity.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(refreshTokenLifetime) * time.Hour * 24),
	}

	err := repo.db.Model(&entity.User{ID: userID}).Association("RefreshTokens").Append(&refreshTokenToCreate)
	return err
}

func (repo *PostgresUserRepository) Activate(userID uint) error {
	user, err := repo.ByID(userID)
	if err != nil {
		return err
	}

	user.Active = true

	result := repo.db.Save(&user)
	return result.Error
}

func (repo *PostgresUserRepository) ChangeNickname(userID uint, newUsername string) error {
	user, err := repo.ByID(userID)
	if err != nil {
		return err
	}

	user.MinecraftCredential.Username = newUsername

	result := repo.db.Save(&user.MinecraftCredential)
	return result.Error
}
