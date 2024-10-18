package postgres

import (
	"backend/internal/entity"
	"backend/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresGameRepository struct {
	db *gorm.DB
}

func NewPostgresGameRepository(db *gorm.DB) *PostgresGameRepository {
	return &PostgresGameRepository{db: db}
}

func (repo *PostgresGameRepository) GenerateAccessTokenForUserAndSave(userID uint) (string, error) {
	timestampString := fmt.Sprint(time.Now().Unix())
	accessToken := utils.MD5ToString(timestampString)

	var user entity.User
	result := repo.db.Preload(clause.Associations).Where(&entity.User{ID: userID}).First(&user)
	if result.Error != nil {
		return accessToken, result.Error
	}

	user.MinecraftCredential.AccessToken = accessToken

	repo.db.Save(&user.MinecraftCredential)
	return accessToken, nil
}

func (repo *PostgresGameRepository) UserByAccessTokenAndUUID(accessToken, uuid string) (entity.User, error) {
	var user entity.User
	var minecraftCredential entity.MinecraftCredential

	result := repo.db.Preload(clause.Associations).Where("uuid = ?", uuid).Where("access_token = ?", accessToken).First(&minecraftCredential)
	if result.Error != nil {
		return user, result.Error
	}

	result = repo.db.Preload(clause.Associations).Where("minecraft_credential_id = ?", minecraftCredential.ID).First(&user)

	return user, result.Error
}

func (repo *PostgresGameRepository) UserByUsernameAndServerID(username, serverID string) (entity.User, error) {
	var user entity.User
	var minecraftCredential entity.MinecraftCredential

	result := repo.db.Preload(clause.Associations).Where("username = ?", username).Where("server_id = ?", serverID).First(&minecraftCredential)
	if result.Error != nil {
		return user, result.Error
	}

	result = repo.db.Preload(clause.Associations).Where("minecraft_credential_id = ?", minecraftCredential.ID).First(&user)
	return user, result.Error
}

func (repo *PostgresGameRepository) UserByUUID(uuid string) (entity.User, error) {
	var user entity.User
	var minecraftCredential entity.MinecraftCredential

	result := repo.db.Preload(clause.Associations).Where("uuid = ?", uuid).First(&minecraftCredential)
	if result.Error != nil {
		return user, result.Error
	}

	result = repo.db.Preload(clause.Associations).Where("minecraft_credential_id = ?", minecraftCredential.ID).First(&user)

	return user, result.Error
}

func (repo *PostgresGameRepository) ApplyServerIDToUser(userID uint, serverID string) error {
	var user entity.User
	result := repo.db.Preload(clause.Associations).Where(&entity.User{ID: userID}).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.MinecraftCredential.ServerID = serverID

	repo.db.Save(&user.MinecraftCredential)
	return nil
}
