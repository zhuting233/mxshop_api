package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		user := claims.(*models.CustomClaims)

		if user.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
