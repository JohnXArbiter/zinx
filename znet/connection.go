package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	Conn     *net.TCPConn   // 当前连接的socket TCP套接字
	ConnID   uint32         // 连接的ID
	isClosed bool           // 当前的连接状态
	ExitChan chan bool      // 告诉当前连接已经退出/停止 channel
	Router   ziface.IRouter // 当前连接处理的方法router
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}
		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			data: buf,
		}
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
		// 从路由中，找到注册绑定的Conn对应的router调用
	}
}

// Start 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID= ", c.ConnID)
	// 启动从当前连接的读数据业务
	go c.StartReader()
	// TODO
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	c.Conn.Close()
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID

}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
