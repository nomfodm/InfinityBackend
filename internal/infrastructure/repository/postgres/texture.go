package postgres

import (
	"backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresTextureRepository struct {
	db *gorm.DB
}

func NewPostgresTextureRepository(db *gorm.DB) *PostgresTextureRepository {
	return &PostgresTextureRepository{db: db}
}

func (repo *PostgresTextureRepository) SkinByID(skinID uint) (entity.Skin, error) {
	var skin entity.Skin
	result := repo.db.Preload(clause.Associations).Where(&entity.Skin{ID: skinID}).First(&skin)
	return skin, result.Error
}

func (repo *PostgresTextureRepository) CapeByID(capeID uint) (entity.Cape, error) {
	var cape entity.Cape
	result := repo.db.Preload(clause.Associations).Where(&entity.Cape{ID: capeID}).First(&cape)
	return cape, result.Error
}

func (repo *PostgresTextureRepository) SkinByHash(hash string) (entity.Skin, error) {
	var skin entity.Skin
	result := repo.db.Preload(clause.Associations).Where("texture_hash = ?", hash).First(&skin)
	return skin, result.Error
}

func (repo *PostgresTextureRepository) CapeByHash(hash string) (entity.Cape, error) {
	var cape entity.Cape
	result := repo.db.Preload(clause.Associations).Where("texture_hash = ?", hash).First(&cape)
	return cape, result.Error
}

func (repo *PostgresTextureRepository) CreateSkin(hash string) (uint, error) {
	skinToCreate := entity.Skin{
		TextureHash: hash,
	}
	result := repo.db.Create(&skinToCreate)
	return skinToCreate.ID, result.Error
}

func (repo *PostgresTextureRepository) CreateCape(hash string) (uint, error) {
	capeToCreate := entity.Cape{
		TextureHash: hash,
	}
	result := repo.db.Create(&capeToCreate)
	return capeToCreate.ID, result.Error
}

func (repo *PostgresTextureRepository) SetSkinToUser(userID, skinID uint) error {
	var user entity.User
	result := repo.db.Where(&entity.User{ID: userID}).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.SkinID = &skinID
	repo.db.Save(&user)
	return nil
}

func (repo *PostgresTextureRepository) SetCapeToUser(userID, capeID uint) error {
	var user entity.User
	result := repo.db.Where(&entity.User{ID: userID}).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.CapeID = &capeID
	repo.db.Save(&user)
	return nil
}
