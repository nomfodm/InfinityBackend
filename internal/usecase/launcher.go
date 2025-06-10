package usecase

import (
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
	"time"
)

type LauncherUseCaseImpl struct {
	repo repository.LauncherRepository
}

func NewLauncherUseCaseImpl(repo repository.LauncherRepository) *LauncherUseCaseImpl {
	return &LauncherUseCaseImpl{repo: repo}
}

func (uc *LauncherUseCaseImpl) ActualLauncherVersion() (entity.LauncherVersion, error) {
	return uc.repo.GetLatestLauncherVersion()
}

func (uc *LauncherUseCaseImpl) RegisterNewUpdate(version, sha256, downloadUrl string) (entity.LauncherVersion, error) {
	newVersion := entity.LauncherVersion{
		DownloadUrl: downloadUrl,
		SHA256:      sha256,
		Version:     version,
		ReleaseDate: time.Now(),
	}
	err := uc.repo.CreateNewLauncherVersion(newVersion)
	return newVersion, err
}
