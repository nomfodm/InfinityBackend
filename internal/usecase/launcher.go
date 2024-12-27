package usecase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/infrastructure/repository"
	"io"
	"net/http"
	"os"
)

type LauncherUseCaseImpl struct {
	repo repository.LauncherRepository
}

func NewLauncherUseCaseImpl(repo repository.LauncherRepository) *LauncherUseCaseImpl {
	return &LauncherUseCaseImpl{repo: repo}
}

func (uc *LauncherUseCaseImpl) CheckForUpdates(clientVersion, clientHash string) (actualVersion, actualHash string, err error) {
	lastVersion, err := uc.repo.GetCurrentVersionInformation()
	if err != nil {
		return "", "", err
	}

	return lastVersion.CurrentVersion, lastVersion.CurrentBinarySHA256, nil
}

func (uc *LauncherUseCaseImpl) CheckForANewUpdate() error {
	ghRepoUrlAPI := os.Getenv("GITHUB_LAUNCHER_REPOSITORY_API_URL")
	response, err := http.Get(ghRepoUrlAPI + "/releases/latest")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	type releaseResponse struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
		}
	}

	responseBody, _ := io.ReadAll(response.Body)

	var responseBodyJson releaseResponse
	err = json.Unmarshal(responseBody, &responseBodyJson)
	if err != nil {
		return err
	}

	lastVersion, err := uc.repo.GetCurrentVersionInformation()
	if err != nil {
		lastVersion = entity.LauncherVersion{}
	}

	hashResponse, err := http.Get(responseBodyJson.Assets[0].BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer hashResponse.Body.Close()
	if hashResponse.StatusCode != 200 { // no releases at all
		return err
	}

	launcherUpdateBinaryData, err := io.ReadAll(hashResponse.Body)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	if _, err = io.Copy(hasher, bytes.NewBuffer(launcherUpdateBinaryData)); err != nil {
		return err
	}

	ghVersion := responseBodyJson.TagName
	ghHash := hex.EncodeToString(hasher.Sum(nil))

	if lastVersion.CurrentVersion != ghVersion {
		if lastVersion.CurrentBinarySHA256 == ghHash {
			err = uc.repo.ModifyExistingVersion(entity.LauncherVersion{
				ID:                  lastVersion.ID,
				CurrentVersion:      ghVersion,
				CurrentBinarySHA256: ghHash,
			})
			return nil
		}
		err = uc.repo.ReleaseNewVersion(entity.LauncherVersion{
			CurrentVersion:      ghVersion,
			CurrentBinarySHA256: ghHash,
		})

		return nil
	}
	if lastVersion.CurrentBinarySHA256 != ghHash {
		err = uc.repo.ModifyExistingVersion(entity.LauncherVersion{
			ID:                  lastVersion.ID,
			CurrentVersion:      ghVersion,
			CurrentBinarySHA256: ghHash,
		})
	}

	return nil
}
