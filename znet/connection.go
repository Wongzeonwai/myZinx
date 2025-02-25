package znet

import (
	"fmt"
	"go-zinx/ziface"
	"net"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID
	ConnID uint32
	// 当前连接状态
	IsClosed bool
	// 告知已经停止的channel
	ExitChan chan bool
	// 该连接处理的方法Router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		IsClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			data: buf,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		c.Router.PreHandle(req)
		c.Router.Handle(req)
		c.Router.PostHandle(req)
	}
}

func (c *Connection) Start() {
	fmt.Println("conn start... connID:", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true

	c.Conn.Close()

	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
