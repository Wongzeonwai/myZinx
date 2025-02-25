package main

import (
	"fmt"
	"go-zinx/ziface"
	"go-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("PingRouter PreHandle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping")); err != nil {
		fmt.Println("PingRouter PreHandle err:", err)
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("ping")); err != nil {
		fmt.Println("PingRouter Handle err:", err)
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("PingRouter PostHandle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping")); err != nil {
		fmt.Println("PingRouter PostHandle err:", err)
	}
}

func main() {
	s := znet.NewServer("[ZINXV0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
