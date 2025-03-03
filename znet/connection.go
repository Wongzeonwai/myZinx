package znet

import (
	"errors"
	"fmt"
	"go-zinx/ziface"
	"io"
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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read err:", err)
		//	continue
		//}
		// 创建一个拆包对象
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read head data error:", err)
			break
		}

		msgData, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack data error:", err)
			break
		}

		var data []byte
		if msgData.GetMsgLen() > 0 {
			data = make([]byte, msgData.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msgData.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msgData,
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

func (c *Connection) Send(msgID uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection is closed")
	}
	// 封包 msglen/msgid/msgdata
	dp := NewDataPack()
	msg := NewMessage(msgID, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		return err
	}
	return nil
}
