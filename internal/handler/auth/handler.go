package auth

import (
	"backend/internal/usecase"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func jsonError(ctx *gin.Context, code int, err string, errDetail error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error":  err,
		"detail": errDetail.Error(),
	})
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var request signUpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonError(ctx, 400, "Validation error", err)
		return
	}

	username := request.Username
	email := request.Email
	password := request.Password

	err := h.uc.SignUp(username, email, password)
	if err != nil {
		jsonError(ctx, 400, "User registration error", err)
		return
	}

	ctx.JSON(201, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var request signInRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonError(ctx, 400, "Validation error", err)
		return
	}

	username := request.Username
	password := request.Password

	accessToken, refreshToken, err := h.uc.SignIn(username, password)
	if err != nil {
		jsonError(ctx, 400, "User authentication error", err)
		return
	}

	refreshTokenLifetime, _ := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_LIFETIME"), 10, 64)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("infinitymcRefreshToken", refreshToken, int(refreshTokenLifetime)*24*60*60, "/", "localhost:8000", true, true)

	ctx.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("infinitymcRefreshToken")
	refreshToken := cookie

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("infinitymcRefreshToken", "", 1, "/", "/", false, false)

	if err != nil {
		var request logoutRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			jsonError(ctx, 400, "Validation error", err)
			return
		}
		refreshToken = request.RefreshToken
	}

	h.uc.Logout(refreshToken)

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("infinitymcRefreshToken")
	refreshToken := cookie
	requestedFromBrowser := true

	if err != nil {
		var request refreshRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			jsonError(ctx, 400, "Validation error", err)
			return
		}
		refreshToken = request.RefreshToken
		requestedFromBrowser = false
	}

	accessToken, refreshToken, err := h.uc.Refresh(refreshToken)
	if err != nil {
		jsonError(ctx, 400, "Token renewing error", err)
		return
	}

	refreshTokenLifetime, _ := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_LIFETIME"), 10, 64)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("infinitymcRefreshToken", refreshToken, int(refreshTokenLifetime)*24*60*60, "/", "localhost:8000", true, true)

	if !requestedFromBrowser {
		ctx.JSON(200, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"accessToken": accessToken,
	})
}

func (h *AuthHandler) Activate(ctx *gin.Context) {
	activationCodeParam := ctx.Query("code")

	activationCode := strings.Split(activationCodeParam, "/")
	if len(activationCode) != 2 {
		jsonError(ctx, 400, "Validation error", errors.New("activation code must be uuid/uint type"))
		return
	}
	uuid := activationCode[0]
	userID, err := strconv.ParseUint(activationCode[1], 10, 32)
	if err != nil {
		jsonError(ctx, 400, "Validation error", errors.New("activation code must be uuid/uint type"))
		return
	}

	err = h.uc.Activate(uint(userID), uuid)
	if err != nil {
		jsonError(ctx, 400, "Activation error", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
