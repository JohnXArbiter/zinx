package znet

import (
	"fmt"
	"net"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name       string             // 服务器的名称
	IPVersion  string             // 服务器绑定的ip版本
	IP         string             // 服务器监听的ip
	Port       int                // 服务器监听的端口
	MsgHandler ziface.IMsgHandler // 当前server的消息管理模块，用来绑定和对应的处理业务API关系
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启workerPool
		s.MsgHandler.StartWorkerPool()

		// 1 获取一个TCP的Addr func ResolveTCPAddr(network string, address string) (*TCPAddr, error)
		addr, err := net.ResolveTCPAddr(s.IPVersion, s.IP+":"+strconv.FormatInt(int64(s.Port), 10))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址 func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, " success, Listening...")

		var cid uint32
		cid = 0

		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP() // (*TCPConn, error)
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			// 客户端已经建立连接，开始做业务

			// 将处理新连接的业务方法 和 conn 进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()
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

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Success")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IService {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
}
