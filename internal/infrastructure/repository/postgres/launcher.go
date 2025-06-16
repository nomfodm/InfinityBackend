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

func (repo *PostgresLauncherRepository) LatestLauncherVersion() (entity.LauncherVersion, error) {
	var launcherVersion entity.LauncherVersion
	result := repo.db.Last(&launcherVersion)
	return launcherVersion, result.Error
}

func (repo *PostgresLauncherRepository) CreateNewLauncherVersion(version entity.LauncherVersion) error {
	result := repo.db.Create(&version)
	return result.Error
}

func (repo *PostgresLauncherRepository) LastMandatoryVersion() (entity.LauncherVersion, error) {
	var launcherVersion entity.LauncherVersion
	result := repo.db.Last(&launcherVersion, "mandatory = true")
	return launcherVersion, result.Error
}
