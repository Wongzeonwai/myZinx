package ziface

/*
	Request接口：
		实际上是把从客户端请求的连接信息和数据封装到一个Request中
*/

type IRequest interface {
	// 得到连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
