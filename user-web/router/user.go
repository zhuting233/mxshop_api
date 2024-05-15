package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	zap.S().Info("配置用户相关url")
	{
		userRouter.GET("list", api.GetUserList)
		//userRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userRouter.POST("pwd_login", api.PassWordLogin)
		userRouter.POST("register", api.Register)
	}
}
