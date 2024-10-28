package user

import (
	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
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

func jsonError(ctx *gin.Context, code int, err string, errDetail error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error":  err,
		"detail": errDetail.Error(),
	})
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
		jsonError(ctx, 400, "File upload error", err)
		return
	}

	hash, err := h.uc.UploadSkin(user, *skinFileHeader)
	if err != nil {
		jsonError(ctx, 500, "Skin uploading error "+hash, err)
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
		jsonError(ctx, 400, "File upload error", err)
		return
	}

	hash, err := h.uc.UploadCape(user, *capeFileHeader)
	if err != nil {
		jsonError(ctx, 500, "Cape uploading error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"user":     user.ID,
		"filename": capeFileHeader.Filename,
		"filesize": capeFileHeader.Size,
		"capeHash": hash,
	})

}
