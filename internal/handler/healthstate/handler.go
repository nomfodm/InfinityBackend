package healthstate

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"github.com/nomfodm/InfinityBackend/internal/utils"
)

type HealthStateHandler struct {
	uc usecase.HealthStateUseCase
}

func NewHealthStateHandler(uc usecase.HealthStateUseCase) *HealthStateHandler {
	return &HealthStateHandler{uc: uc}
}

func (h *HealthStateHandler) Index(ctx *gin.Context) {
	status := GetStatusFromCtx(ctx)
	ctx.JSON(200, gin.H{
		"status": entity.ServerStatuses[status.Status],
		"time":   time.Now(),
	})
}

func (h *HealthStateHandler) SetStatus(ctx *gin.Context) {
	newStatus, _ := ctx.GetQuery("newStatus")
	if newStatus == "" {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error":  "Enter server status",
			"detail": "",
		})
		return
	}

	newStatusInt, err := strconv.Atoi(newStatus)
	if err != nil {
		utils.JsonError(ctx, 500, "String to integer conversion error", err)
		return
	}
	err = h.uc.SetHealthState(newStatusInt)
	if err != nil {
		utils.JsonError(ctx, 500, "Health state change error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": entity.ServerStatuses[newStatusInt],
		"time":   time.Now(),
	})
}
