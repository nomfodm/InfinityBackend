package status

import (
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"net/http"
	"time"
)

type StatusMiddleware struct {
	uc usecase.ServerStatusUseCase
}

func NewStatusMiddleware(uc usecase.ServerStatusUseCase) gin.HandlerFunc {
	return (&StatusMiddleware{uc}).Handle
}

func statusError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

func StatusNotOK(ctx *gin.Context, status int) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"status": entity.ServerStatuses[status],
		"time":   time.Now(),
	})
}

func GetStatusFromCtx(ctx *gin.Context) entity.ServerStatus {
	serverStatusRaw, _ := ctx.Get("serverstatus")
	return serverStatusRaw.(entity.ServerStatus)
}

func (m *StatusMiddleware) Handle(ctx *gin.Context) {
	status, err := m.uc.CurrentServerStatus()
	if err != nil {
		statusError(ctx, err)
		return
	}

	if status.Status != entity.ServerStatusWorking {
		StatusNotOK(ctx, status.Status)
	}

	ctx.Set("serverstatus", status)
}
