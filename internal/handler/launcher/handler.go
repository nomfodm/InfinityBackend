package launcher

import (
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"github.com/nomfodm/InfinityBackend/internal/utils"
	"net/http"
)

type LauncherHandler struct {
	uc usecase.LauncherUseCase
}

func NewLauncherHandler(uc usecase.LauncherUseCase) *LauncherHandler {
	return &LauncherHandler{uc: uc}
}

func (h *LauncherHandler) ActualVersion(ctx *gin.Context) {
	version, err := h.uc.ActualLauncherVersion()
	if err != nil {
		utils.JsonError(ctx, 500, "ActualLauncherVersion error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"version":     version.Version,
		"sha256":      version.SHA256,
		"url":         version.DownloadUrl,
		"releaseDate": version.ReleaseDate,
		"mandatory":   version.Mandatory,
	})
}

func (h *LauncherHandler) DownloadLauncher(ctx *gin.Context) {
	version, err := h.uc.ActualLauncherVersion()
	if err != nil {
		utils.JsonError(ctx, 500, "ActualLauncherVersion error", err)
		return
	}

	ctx.Redirect(http.StatusFound, version.DownloadUrl)
}

func (h *LauncherHandler) RegisterUpdate(ctx *gin.Context) {
	var request registerUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.JsonError(ctx, 400, "Validation error", err)
		return
	}

	newVersion, err := h.uc.RegisterNewUpdate(request.Version, request.SHA256, request.DownloadUrl, request.Mandatory)
	if err != nil {
		utils.JsonError(ctx, 500, "RegisterNewUpdate error", err)
		return
	}

	ctx.JSON(200, newVersion)
}

func (h *LauncherHandler) LastMandatory(ctx *gin.Context) {
	version, err := h.uc.LastMandatoryVersion()
	if err != nil {
		utils.JsonError(ctx, 500, "LastMandatoryVersion error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"version":     version.Version,
		"sha256":      version.SHA256,
		"url":         version.DownloadUrl,
		"releaseDate": version.ReleaseDate,
		"mandatory":   version.Mandatory,
	})
}
