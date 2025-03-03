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
			// 将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将服务器资源、状态或者已经开辟的连接信息停止或回收
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

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
