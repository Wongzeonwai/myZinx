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
	// 当前连接绑定的处理业务方法
	handleAPI ziface.HandleFunc
	// 告知已经停止的channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, handleAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		handleAPI: handleAPI,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", err)
			continue
		}

		// 调用当前连接绑定的handleAPI
		if err := c.handleAPI(c.Conn, buf[:cnt], cnt); err != nil {
			fmt.Println("handle err:", err)
			break
		}
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
