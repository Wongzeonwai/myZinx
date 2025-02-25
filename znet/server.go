package znet

import (
	"fmt"
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
}

func (s *Server) Start() {
	fmt.Printf("server start at：%s，port：%d\n", s.IP, s.Port)
	go func() {
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
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err：", err)
				continue
			}

			// 客户端已经建立连接，做一些业务，做一个最大512字节的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Println("conn read err：", err)
						continue
					}
					fmt.Println(string(buf[:n]))
					if _, err = conn.Write(buf[:n]); err != nil {
						fmt.Println("conn write err：", err)
						continue
					}
				}
			}()
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

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
