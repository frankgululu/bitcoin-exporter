package main

import (
	"bitcoin-exporter/config"
	"bitcoin-exporter/logger"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. 定义监控指标
var bitcoinSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "bitcoin_sync_status",
	Help: "Current sync status of the Bitcoin node",
})

var bitcoinBlockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "bitcoin_block_height",
	Help: "Current block height of the Bitcoin node",
})

// 2. 注册到promehtues
func init() {
	prometheus.MustRegister(bitcoinSyncStatus, bitcoinBlockHeight)
}

// 3. 处理metrics
// 调用bitcoin rpc获取当前节点的状态信息
func UpdateBlockchainMetrics(rpcHost, rpcUser, rpcPass string, rpcPort int, useSSL bool) error {
	// 配置 RPC 连接选项
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", rpcHost, rpcPort),
		User:         rpcUser,
		Pass:         rpcPass,
		HTTPPostMode: true,    // 使用HTTP POST模式
		DisableTLS:   !useSSL, // 根据useSSL决定是否使用TLS
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		logger.Logger.Error("创建RPC client失败:", "error", err)
		return nil
	}

	defer client.Shutdown()

	blockchainInfo, err := client.GetBlockChainInfo()

	if err != nil {
		logger.Logger.Error("获取区块信息失败:", "error", err)
		return nil
	}

	//更新prometheus metrics
	if blockchainInfo.InitialBlockDownload {
		bitcoinSyncStatus.Set(1) //true syncing
	} else {
		bitcoinSyncStatus.Set(0) //false synced
	}

	bitcoinBlockHeight.Set(float64(blockchainInfo.Blocks))

	return nil

}

//4. 启动服务

func main() {
	//读取config包里序列化过来的环境变量的值
	config.NewConfig()
	rpcHost := config.Rpc.Host
	rpcUser := config.Rpc.User
	rpcPass := config.Rpc.Pass
	rpcPort := config.Rpc.Port
	useSSL := config.Rpc.Ssl

	var mutex sync.Mutex
	wait := sync.WaitGroup{}

	// 启动 Goroutine，每隔 10 秒更新一次区块链指标
	go func() {
		for {
			wait.Add(1)
			mutex.Lock()
			err := UpdateBlockchainMetrics(rpcHost, rpcUser, rpcPass, rpcPort, useSSL)
			if err != nil {
				logger.Logger.Error("update metric to promethues err", "err", err)
				//slog.Error("update metric to promethues err", slog.Any("err", err))
				return
			}
			mutex.Unlock()
			time.Sleep((10 * time.Second))

		}

	}()

	addr := "0.0.0.0:2024"

	http.Handle("/metrics", promhttp.Handler())
	logger.Logger.Info("Staring server at", "address", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Logger.Error("Server failed to start:", "err", err)

	}

}
