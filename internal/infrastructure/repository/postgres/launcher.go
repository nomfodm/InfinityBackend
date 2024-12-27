package postgres

import (
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"gorm.io/gorm"
)

type PostgresLauncherRepository struct {
	db *gorm.DB
}

func NewPostgresLauncherRepository(db *gorm.DB) *PostgresLauncherRepository {
	return &PostgresLauncherRepository{db: db}
}

func (repo *PostgresLauncherRepository) GetCurrentVersionInformation() (entity.LauncherVersion, error) {
	var launcherVersion entity.LauncherVersion
	result := repo.db.Last(&launcherVersion)
	return launcherVersion, result.Error
}

func (repo *PostgresLauncherRepository) ReleaseNewVersion(version entity.LauncherVersion) error {
	result := repo.db.Create(&version)
	return result.Error
}

func (repo *PostgresLauncherRepository) ModifyExistingVersion(version entity.LauncherVersion) error {
	result := repo.db.Save(&version)
	return result.Error
}
