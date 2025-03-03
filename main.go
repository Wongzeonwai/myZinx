package main

import (
	"fmt"
	"go-zinx/ziface"
	"go-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	// 先读取客户端数据，再回写ping
	fmt.Println("recv from client:msgID=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().Send(1, []byte("ping"))
	if err != nil {
		fmt.Println("send ping err:", err)
	}
}

func main() {
	s := znet.NewServer("[ZINXV0.5]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
