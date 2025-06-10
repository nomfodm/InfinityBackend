package launcher

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type AdminAccessMiddleware struct {
}

func NewAdminAccessMiddleware() gin.HandlerFunc {
	return (&AdminAccessMiddleware{}).Handle
}

func (m *AdminAccessMiddleware) Handle(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader != os.Getenv("ADMIN_PASSWORD") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":  "Unauthorized",
			"detail": "",
		})
		return
	}
}
