package ziface

import "net"

type IConnection interface {
	// 启动连接
	Start()
	// 停止连接
	Stop()
	// 获取当前连接绑定的socket
	GetTCPConnection() *net.TCPConn
	// 获取当前连接的ID
	GetID() uint32
	// 获取远程客户端的TCP状态，IP和Port
	GetRemoteAddr() net.Addr
	// 发送数据，将数据发送给远程客户端
	Send(data []byte) error
}

// 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
