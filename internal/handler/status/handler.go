package status

import (
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"time"
)

type StatusHandler struct {
	uc usecase.ServerStatusUseCase
}

func NewStatusHandler(uc usecase.ServerStatusUseCase) *StatusHandler {
	return &StatusHandler{uc}
}

func (h *StatusHandler) Index(ctx *gin.Context) {
	status := GetStatusFromCtx(ctx)
	ctx.JSON(200, gin.H{
		"status": entity.ServerStatuses[status.Status],
		"time":   time.Now(),
	})
}

func (h *StatusHandler) SetStatus(ctx *gin.Context) {
	newStatus, _ := ctx.GetQuery("newStatus")
	if newStatus == "" {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error":  "Enter server status",
			"detail": "",
		})
		return
	}
}
