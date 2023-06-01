package znet

import (
	"fmt"
	"net"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name        string                        // 服务器的名称
	IPVersion   string                        // 服务器绑定的ip版本
	IP          string                        // 服务器监听的ip
	Port        int                           // 服务器监听的端口
	MsgHandler  ziface.IMsgHandler            // 当前server的消息管理模块，用来绑定和对应的处理业务API关系
	ConnMgr     ziface.IConnManager           // 该server的链接管理器
	OnConnStart func(conn ziface.IConnection) // 该server创建链接之后自动调用hook函数--OnConnStart
	onConnStop  func(conn ziface.IConnection) // 该server销毁链接之前自动调用的hook函数--OnConnStop
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
		fmt.Println("start Zinx server success,", s.Name, " success, Listening ...")

		var cid uint32
		cid = 0
		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP() // (*TCPConn, error)
			if err != nil {
				fmt.Println("Accept error ", err)
				continue
			}
			// 最大连接个数判断，超过最大个数，不接受新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 将处理新连接的业务方法 和 conn 进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或者回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	// 阻塞状态
	select {}
}

// AddRouter 路由功能:给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化Server模块的方法
func NewServer() ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnMgr(),
	}
}

// SetOnConnStart 注册OnConnStart hook的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop hook的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.onConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart hook的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop hook的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("----> Call OnConnStop() ...")
		s.onConnStop(conn)
	}
}
