package znet

import "go-zinx/ziface"

// 实现路由时，先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

func (br *BaseRouter) Handle(req ziface.IRequest) {}

func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
