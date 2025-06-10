package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"github.com/nomfodm/InfinityBackend/internal/utils"
)

type UserHandler struct {
	uc usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func parseUserFromContext(ctx *gin.Context) entity.User {
	userRaw, _ := ctx.Get("user")
	return userRaw.(entity.User)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	user := parseUserFromContext(ctx)

	skinHash, capeHash := h.uc.SkinCapeHashes(user)
	textures := textures{SkinHash: skinHash, CapeHash: capeHash}

	response := meResponse{ID: user.ID,
		Username:            user.Username,
		Email:               user.Email,
		Active:              user.Active,
		RegisteredAt:        user.RegisteredAt,
		Textures:            textures,
		MinecraftCredential: minecraftCredential{Username: user.MinecraftCredential.Username, UUID: user.MinecraftCredential.UUID},
	}

	ctx.JSON(200, response)
}

func (h *UserHandler) Skin(ctx *gin.Context) {
	user := parseUserFromContext(ctx)

	skinFileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.JsonError(ctx, 400, "File upload error", err)
		return
	}

	hash, err := h.uc.UploadSkin(user, *skinFileHeader)
	if err != nil {
		utils.JsonError(ctx, 500, "Skin uploading error "+hash, err)
		return
	}

	ctx.JSON(200, gin.H{
		"user":     user.ID,
		"filename": skinFileHeader.Filename,
		"filesize": skinFileHeader.Size,
		"skinHash": hash,
	})
}

func (h *UserHandler) Cape(ctx *gin.Context) {
	user := parseUserFromContext(ctx)

	capeFileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.JsonError(ctx, 400, "File upload error", err)
		return
	}

	hash, err := h.uc.UploadCape(user, *capeFileHeader)
	if err != nil {
		utils.JsonError(ctx, 500, "Cape uploading error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"user":     user.ID,
		"filename": capeFileHeader.Filename,
		"filesize": capeFileHeader.Size,
		"capeHash": hash,
	})
}

func (h *UserHandler) Nickname(ctx *gin.Context) {
	user := parseUserFromContext(ctx)

	newNickname := ctx.Query("new_nickname")
	if newNickname == "" {
		utils.JsonError(ctx, 400, "Nickname is required", errors.New(""))
		return
	}

	err := h.uc.ChangeNickname(user, newNickname)

	if err != nil {
		utils.JsonError(ctx, 500, "Nickname change error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
