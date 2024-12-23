package usecase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type LauncherUseCaseImpl struct {
}

func NewLauncherUseCaseImpl() *LauncherUseCaseImpl {
	return &LauncherUseCaseImpl{}
}

func (uc *LauncherUseCaseImpl) CheckForUpdates(clientVersion, clientHash string) (actualVersion, actualHash string, isUpdates bool, err error) {
	ghRepoUrlAPI := os.Getenv("GITHUB_LAUNCHER_REPOSITORY_API_URL")
	response, err := http.Get(ghRepoUrlAPI + "/releases/latest")
	if err != nil {
		return "", "", false, err
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
		return "", "", false, err
	}

	actualVersion = responseBodyJson.TagName
	if actualVersion != clientVersion {
		return actualVersion, "", true, nil
	}

	hashResponse, err := http.Get(responseBodyJson.Assets[0].BrowserDownloadURL)
	if err != nil {
		return "", "", false, err
	}
	defer hashResponse.Body.Close()
	if hashResponse.StatusCode != 200 {
		return "", "", false, nil // no release
	}

	launcherUpdateBinaryData, err := io.ReadAll(hashResponse.Body)
	if err != nil {
		return "", "", false, err
	}

	hasher := sha256.New()
	if _, err = io.Copy(hasher, bytes.NewBuffer(launcherUpdateBinaryData)); err != nil {
		return "", "", false, err
	}

	actualHash = hex.EncodeToString(hasher.Sum(nil))
	if clientHash != actualHash {
		return actualVersion, actualHash, true, nil
	}

	return "", "", false, nil
}
