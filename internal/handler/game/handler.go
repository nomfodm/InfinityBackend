package game

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"github.com/nomfodm/InfinityBackend/internal/utils"
)

type GameHandler struct {
	uc usecase.GameUseCase
}

func NewGameHandler(uc usecase.GameUseCase) *GameHandler {
	return &GameHandler{uc: uc}
}

func parseUserFromContext(ctx *gin.Context) entity.User {
	userRaw, _ := ctx.Get("user")
	return userRaw.(entity.User)
}

func minecraftError(ctx *gin.Context, code int, err string, errDetail error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error":        err,
		"errorMessage": errDetail.Error(),
		"cause":        "",
	})
}

func (h *GameHandler) Launcher(ctx *gin.Context) {
	user := parseUserFromContext(ctx)
	username, uuid, accessToken, err := h.uc.Launcher(user)
	if err != nil {
		utils.JsonError(ctx, 500, "Cant generate user data", err)
		return
	}

	ctx.JSON(200, gin.H{
		"username":    username,
		"uuid":        uuid,
		"accessToken": accessToken,
	})
}

func (h *GameHandler) Join(ctx *gin.Context) {
	var request joinRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		minecraftError(ctx, 400, "Validation error", err)
		return
	}

	accessToken := request.AccessToken
	UUID := request.SelectedProfile
	serverID := request.ServerID

	err := h.uc.Join(accessToken, utils.AddHyphenToUUID(UUID), serverID)
	if err != nil {
		minecraftError(ctx, 400, err.Error(), errors.New("играть на сервере можно только через лаунчер Infinity"))
		return
	}

	ctx.Status(200)
}

func (h *GameHandler) HasJoined(ctx *gin.Context) {
	username := ctx.Query("username")
	serverID := ctx.Query("serverId")
	if username == "" || serverID == "" {
		ctx.Status(400)
		return
	}
	user, err := h.uc.HasJoined(username, serverID)
	if err != nil {
		ctx.Status(400)
		return
	}

	user, skin, cape, err := h.uc.Profile(user.MinecraftCredential.UUID.String())
	if err != nil {
		ctx.Status(400)
		return
	}

	profile, err := Profile(user, skin, cape)
	if err != nil {
		ctx.Status(500)
		return
	}

	ctx.JSON(200, profile)
}

func (h *GameHandler) Profile(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if uuid == "" || !(len(uuid) == 32 || len(uuid) == 36) {
		ctx.Status(400)
		return
	}

	user, skin, cape, err := h.uc.Profile(utils.AddHyphenToUUID(uuid))
	if err != nil {
		ctx.Status(500)
		return
	}

	profile, err := Profile(user, skin, cape)
	if err != nil {
		ctx.Status(502)
		return
	}

	ctx.JSON(200, profile)
}
