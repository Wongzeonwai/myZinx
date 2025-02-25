package ziface

type IRouter interface {
	// 处理连接业务之前的方法
	PreHandle(request IRequest)
	// 处理连接业务的主方法
	Handle(request IRequest)
	// 处理连接业务之后的方法
	PostHandle(request IRequest)
}
