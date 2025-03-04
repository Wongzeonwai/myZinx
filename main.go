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

type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	// 先读取客户端数据，再回写ping
	fmt.Println("recv from client:msgID=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().Send(201, []byte("hello welcome to zinx"))
	if err != nil {
		fmt.Println("send ping err:", err)
	}
}

func DoconnectionBegin(conn ziface.IConnection) {
	fmt.Println("<============ DoconnectionBegin begin ==============>")
	if err := conn.Send(202, []byte("hello zinx")); err != nil {
		fmt.Println("send ping err:", err)
	}
}

func DoconnectionEnd(conn ziface.IConnection) {
	fmt.Println("============ DoconnectionEnd end ==============>")
	fmt.Println("conn id=", conn.GetID(), "is logout")
}

func main() {
	s := znet.NewServer("[ZINXV0.9]")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.SetOnConnStart(DoconnectionBegin)
	s.SetOnConnStop(DoconnectionEnd)
	s.Serve()
}
