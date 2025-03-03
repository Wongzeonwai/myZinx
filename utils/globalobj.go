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

	Version        string
	MaxConn        int
	MaxPackageSize uint32
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
		Name:           "ZinxServerAPP",
		Version:        "V0.7",
		TCPPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalObject.Reload()
}
