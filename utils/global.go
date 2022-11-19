package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"server/iface"
)

/*
	存储全局参数，供其他模块使用
	一些参数是可以由用户通过json配置
*/

type GlobalConfig struct {
	TcpServer      iface.IServer
	Name           string `json:"name"`
	Host           string `json:"host"`
	TcpPort        int    `json:"tcpPort"`
	Version        string `json:"version"`
	MaxConn        int    `json:"maxConn"`        // 最大连接数
	MaxPackingSize uint32 `json:"maxPackingSize"` // 当前服务器一次数据包的最大值
}

var GlobalObj *GlobalConfig

// init 初始化对象
func init() {
	// 默认配置
	GlobalObj = &GlobalConfig{
		Name:           "default_server",
		Version:        "1.0",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        10,
		MaxPackingSize: 4096,
	}
	// read json config
	GlobalObj.loadFormJson()
}

func (g *GlobalConfig) loadFormJson() {
	data, err := os.ReadFile("./server_config.json")
	if err != nil {
		fmt.Printf("[Server Reading Config ERROR] Reading config error:%s\n", err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		fmt.Printf("[Server Reading Config ERROR] Reading config error:%s\n", err)
	}
}
