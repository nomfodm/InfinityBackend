package auth

import (
	"github.com/nomfodm/InfinityBackend/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	uc usecase.AuthUseCase
}

func NewAuthMiddleware(uc usecase.AuthUseCase) gin.HandlerFunc {
	return (&AuthMiddleware{uc: uc}).Handle
}

func notAuthorized(ctx *gin.Context, err string, detail string) {
	ctx.AbortWithStatusJSON(401, gin.H{
		"error":  err,
		"detail": detail,
	})
}

func (m *AuthMiddleware) Handle(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		notAuthorized(ctx, "Unauthorized", "provide access token string")
		return
	}

	headerSplitted := strings.Split(authorizationHeader, " ")
	if len(headerSplitted) != 2 {
		notAuthorized(ctx, "Unauthorized", "provide corrent access token string")
		return
	}

	if headerSplitted[0] != "Bearer" {
		notAuthorized(ctx, "Unauthorized", "access token type must be bearer")
		return
	}

	accessToken := headerSplitted[1]

	if accessToken == "null" {
		notAuthorized(ctx, "Unauthorized", "provide corrent access token string")
		return
	}

	user, err := m.uc.RetrieveUserWithAccessToken(accessToken)
	if err != nil {
		notAuthorized(ctx, "Unauthorized", err.Error())
		return
	}

	ctx.Set("user", user)
}
