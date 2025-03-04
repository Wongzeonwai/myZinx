package znet

import (
	"errors"
	"fmt"
	"go-zinx/utils"
	"go-zinx/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前连接属于哪一个server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID
	ConnID uint32
	// 当前连接状态
	IsClosed bool
	// 告知已经停止的channel
	ExitChan chan bool
	// 无缓冲管道，用于读写协程直接的消息同学
	MsgChan    chan []byte
	MsgHandler ziface.IMsgHandle
	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msghandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		MsgHandler: msghandler,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().AddConn(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("start reader ...")
	defer fmt.Println(c.GetRemoteAddr().String(), "reader exi")
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// 专门发送消息给客户端的模块
func (c *Connection) StartWriter() {
	fmt.Println("start writer ...")
	defer fmt.Println(c.GetRemoteAddr().String(), "conn writer exit!")

	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("write data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader退出，此时writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("conn start... connID:", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()
	// 调用连接后的钩子
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	// 调用销毁前的钩子
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true
	c.TcpServer.GetConnManager().RemoveConn(c)
	close(c.ExitChan)
	close(c.MsgChan)
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
	if c.IsClosed == true {
		return errors.New("connection is closed")
	}
	// 封包 msglen/msgid/msgdata
	dp := NewDataPack()
	msg := NewMessage(msgID, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	c.MsgChan <- binaryMsg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if valuem, ok := c.property[key]; ok {
		return valuem, nil
	} else {
		return nil, errors.New("property not found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
