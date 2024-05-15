package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	myvalidator "mxshop-api/user-web/validator"
)

func main() {

	//初始化logger
	initialize.InitLogger()

	//初始化Routers
	Router := initialize.Routers()

	initialize.InitConfig()

	if initialize.InitTrans("zh") != nil {
		zap.S().Panic("InitTrans fail!")
	}

	initialize.InitSrvConn()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Value().(string))
			return t
		})
	}

	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
