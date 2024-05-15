package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/utils"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {

	global.ServerConfig = &config.ServerConfig{}
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "user-web/config.pro.yaml"
	if debug {
		configFileName = "user-web/config.dev.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	zap.S().Infof("使用配置文件: %v", configFileName)
	zap.S().Infof("配置信息: %v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件变化: %v", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息: %v", global.ServerConfig)
	})

}
