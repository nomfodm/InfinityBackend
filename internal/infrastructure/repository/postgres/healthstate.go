package postgres

import (
	"errors"

	"github.com/nomfodm/InfinityBackend/internal/entity"
	"gorm.io/gorm"
)

type PostgresHealthStateRepository struct {
	db *gorm.DB
}

func NewPostgresHealthStateRepository(db *gorm.DB) *PostgresHealthStateRepository {
	return &PostgresHealthStateRepository{db: db}
}

func (repo *PostgresHealthStateRepository) InitHealthState() error {
	var healthState entity.HealthState
	result := repo.db.First(&healthState)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		hS := entity.HealthState{
			Status: entity.ServerStatusWorking,
		}
		result := repo.db.Create(&hS)
		return result.Error
	}
	return nil
}

func (repo *PostgresHealthStateRepository) CurrentHealthState() (entity.HealthState, error) {
	var healthState entity.HealthState
	result := repo.db.First(&healthState)
	return healthState, result.Error
}

func (repo *PostgresHealthStateRepository) SetHealthState(newStatus int) error {
	var healthState entity.HealthState
	result := repo.db.First(&healthState)
	if result.Error != nil {
		return result.Error
	}

	healthState.Status = newStatus
	result = repo.db.Save(&healthState)
	return result.Error
}
