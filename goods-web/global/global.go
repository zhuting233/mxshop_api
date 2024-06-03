package global

import (
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"
	"strings"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig    *config.NacosConfig
	ServerConfig   *config.ServerConfig
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
