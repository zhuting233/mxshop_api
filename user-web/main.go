package main

import (
	"fmt"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils/register/consul"
	myvalidator "mxshop-api/user-web/validator"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {

	//初始化logger
	initialize.InitLogger()
	defer zap.S().Sync()

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
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewV4().String()
	register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)

	go func(){
		zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
		err := Router.Run(fmt.Sprintf("0.0.0.0:%d", global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()

	quit := make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT,syscall.SIGTERM)
	<- quit
	
	if err := register_client.DeRegister(serviceId);err!=nil{
		zap.S().Info("注销失败",err.Error())
	}else{
		zap.S().Info("注销成功")
	}
}
