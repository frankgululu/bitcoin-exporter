package config

import (
	"bitcoin-exporter/logger"

	"github.com/spf13/viper"
)

// 定义结构体
type RpcConfig struct {
	Host string
	User string
	Pass string
	Port int
	Ssl  bool
}

// 实例化结构体

var Rpc RpcConfig

func NewConfig() {
	logger.Logger.Info("init configureation")
	//fmt.Println("init config...")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 启用自动读取环境变量
	//viper.AutomaticEnv()

	// 绑定环境变量
	viper.BindEnv("host", "RPC_HOST")
	viper.BindEnv("user", "RPC_USER")
	viper.BindEnv("pass", "RPC_PASS")
	viper.BindEnv("port", "RPC_PORT")
	viper.BindEnv("ssl", "RPC_SSL")

	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error("Read config file error", "error", err)
		return
		//panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Rpc); err != nil {
		logger.Logger.Error("error unmarshal confilg file,", "error", err)
		//fmt.Println(fmt.Errorf("error unmarshal confilg file,%w", err))
	}
}
