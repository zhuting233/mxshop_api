package main

import (
	"fmt"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
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

	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewV4().String()
	register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name,global.ServerConfig.Tags , serviceId)
	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	go func(){
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
