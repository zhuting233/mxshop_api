package initialize

import (
	"mxshop-api/goods-web/middlewares"
	"mxshop-api/goods-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})
	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)

	return Router
}
