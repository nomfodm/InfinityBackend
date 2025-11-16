package usecase

import (
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
)

type HealthStateUseCaseImpl struct {
	repo repository.HealthStateRepository
}

func NewHealthStateUseCaseImpl(repo repository.HealthStateRepository) *HealthStateUseCaseImpl {
	return &HealthStateUseCaseImpl{repo: repo}
}

func (uc *HealthStateUseCaseImpl) InitHealthState() error {
	return uc.repo.InitHealthState()
}

func (uc *HealthStateUseCaseImpl) CurrentHealthState() (entity.HealthState, error) {
	return uc.repo.CurrentHealthState()
}

func (uc *HealthStateUseCaseImpl) SetHealthState(newStatus int) error {
	return uc.repo.SetHealthState(newStatus)
}
