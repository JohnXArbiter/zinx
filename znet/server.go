package znet

import (
	"fmt"
	"net"
	"strconv"
	"zinx/ziface"
)

type Server struct {
	Name      string // 服务器的名称
	IPVersion string // 服务器绑定的ip版本
	IP        string // 服务器监听的ip
	Port      int    // 服务器监听的端口
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	go func() {
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, s.IP+":"+strconv.FormatInt(int64(s.Port), 10))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, " success, Listening...")
		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 客户端已经建立连接，开始做业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err ", err)
						continue
					}
					fmt.Printf("recv client buf %s，cnt %d\n", buf, cnt)

					// 回显
					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	// TODO
}
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 启动之后做的额外业务

	// 阻塞状态

	select {}
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IService {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9000,
	}
	return s
}
