package launcher

import (
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
)

type LauncherHandler struct {
	uc usecase.LauncherUseCase
}

func NewLauncherHandler(uc usecase.LauncherUseCase) *LauncherHandler {
	return &LauncherHandler{uc: uc}
}

func jsonError(ctx *gin.Context, code int, err string, errDetail error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error":  err,
		"detail": errDetail.Error(),
	})
}

func (h *LauncherHandler) Updates(ctx *gin.Context) {
	actualVersion, actualHash, err := h.uc.CheckForUpdates()
	if err != nil {
		jsonError(ctx, 500, "CheckForUpdates error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"version": actualVersion,
		"hash":    actualHash,
	})
}

func (h *LauncherHandler) CheckForANewUpdate(ctx *gin.Context) {
	err := h.uc.CheckForANewUpdate()
	if err != nil {
		jsonError(ctx, 500, "CheckForANewUpdate error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
