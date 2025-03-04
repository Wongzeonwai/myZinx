package utils

import (
	"encoding/json"
	"go-zinx/ziface"
	"os"
)

type GlobalObj struct {
	TCPServer ziface.IServer
	Host      string
	TCPPort   int
	Name      string

	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32 // 当前业务工作worker池的goroutine数量
	MaxWorkerTaskLen uint32 // zinx框架允许用户最多开辟多少个worker
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerAPP",
		Version:          "V0.9",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          2,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	GlobalObject.Reload()
}
