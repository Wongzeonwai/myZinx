package znet

import (
	"fmt"
	"go-zinx/utils"
	"go-zinx/ziface"
	"net"
)

// 定义一个Server服务类
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前server的消息管理模块
	MsgHandler ziface.IMsgHandle
	// 该server的连接管理器connmanager
	ConnMgr ziface.IConnManager
	// 该server创建连接后自动调用的hook
	OnConnStart func(conn ziface.IConnection)
	// 该server销毁前自动调用的hook
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("server start at：%s，port：%d\n", s.IP, s.Port)
	go func() {
		// 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err：", err)
			return
		}
		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp err：", err)
			return
		}
		fmt.Println("server start success")
		// 阻塞的等待客户端连接，处理客户端连接业务
		var cid uint32
		cid = 0
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err：", err)
				continue
			}
			// 设置是否超出最大连接数量
			if s.ConnMgr.GetConnLen() >= utils.GlobalObject.MaxConn {
				fmt.Println("=======================> too many conn!")
				conn.Close()
				continue
			}
			// 将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将服务器资源、状态或者已经开辟的连接信息停止或回收
	s.ConnMgr.ClearConn()
	fmt.Println(s.Name, " server stop success")
}

func (s *Server) Serve() {
	s.Start()
	// 启动服务器后的额外业务
	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, r ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, r)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册OnConnStart的方法
func (s *Server) SetOnConnStart(hook func(conn ziface.IConnection)) {
	s.OnConnStart = hook
}

// 注册OnConnStop的方法
func (s *Server) SetOnConnStop(hook func(conn ziface.IConnection)) {
	s.OnConnStop = hook
}

// 调用OnConnStart的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

// 调用OnConnStart的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
