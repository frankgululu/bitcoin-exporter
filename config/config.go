package config

import (
	"fmt"

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

	fmt.Println("init config...")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.BindEnv("rpc.host", "RPC_HOST")
	viper.BindEnv("rpc.user", "RPC_USER")
	viper.BindEnv("rpc.pass", "RPC_PASS")
	viper.BindEnv("rpc.port", "RPC_PORT")
	viper.BindEnv("rpc.ssl", "RPC_SSL")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Rpc); err != nil {
		fmt.Println(fmt.Errorf("error unmarshal confile file,%w", err))
	}

}
