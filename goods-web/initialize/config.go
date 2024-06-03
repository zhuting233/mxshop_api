package initialize

import (
	"encoding/json"
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	global.NacosConfig = &config.NacosConfig{}
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "goods-web/config.pro.yaml"
	if debug {
		configFileName = "goods-web/config.dev.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息: %v", global.NacosConfig)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			ContextPath: "/nacos",
			Port:        uint64(global.NacosConfig.Port),
			Scheme:      "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}

	//fmt.Println(content)

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败:%s", err.Error())
	}
	//fmt.Println(serverConfig)
	zap.S().Info(global.ServerConfig)

}

func InitConfig2() {

	global.ServerConfig = &config.ServerConfig{}
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "goods-web/config.pro.yaml"
	if debug {
		configFileName = "goods-web/config.dev.yaml"
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
