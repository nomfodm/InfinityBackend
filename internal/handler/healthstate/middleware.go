package healthstate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
)

type HealthStateMiddleware struct {
	uc usecase.HealthStateUseCase
}

func NewHealthStateMiddleware(uc usecase.HealthStateUseCase) gin.HandlerFunc {
	return (&HealthStateMiddleware{uc: uc}).Handle
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

func GetStatusFromCtx(ctx *gin.Context) entity.HealthState {
	serverStatusRaw, _ := ctx.Get("serverstatus")
	return serverStatusRaw.(entity.HealthState)
}

func (m *HealthStateMiddleware) Handle(ctx *gin.Context) {
	status, err := m.uc.CurrentHealthState()
	if err != nil {
		statusError(ctx, err)
		return
	}

	if status.Status != entity.ServerStatusWorking {
		StatusNotOK(ctx, status.Status)
	}

	ctx.Set("serverstatus", status)
}
