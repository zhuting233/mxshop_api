package router

import (
	"mxshop-api/goods-web/api/goods"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	goodsRouter := Router.Group("goods")
	zap.S().Info("配置商品相关url")
	{
		goodsRouter.GET("", goods.List)

	}
}
