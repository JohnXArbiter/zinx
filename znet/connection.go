package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	TcpServer    ziface.IServer         // 当前connection属于哪个server
	Conn         *net.TCPConn           // 当前连接的socket TCP套接字
	ConnId       uint32                 // 连接的ID
	isClosed     bool                   // 当前的连接状态
	ExitChan     chan struct{}          // Reader告诉Writer去退出 channel
	msgChan      chan []byte            // 无缓冲管道，用于读、写goroutine之间的消息通信
	MsgHandler   ziface.IMsgHandler     // 消息的管理MsgID和对应的处理业务API关系
	property     map[string]interface{} // 链接的属性集合
	propertyLock sync.RWMutex           // 保护链接属性的锁
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnId:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan struct{}, 1),
		property:   make(map[string]interface{}),
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// Start 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnId = ", c.ConnId)
	// 启动从当前连接的读数据业务
	go c.StartReader()
	// 启动从当前连接的写数据业务
	go c.StartWriter()
	// 调用Start后的hook
	c.TcpServer.CallOnConnStart(c)
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running!]")
	defer fmt.Println("connID = ", c.ConnId, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的 msg head 8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error", err)
		}

		// 拆包，得到 msgId 和 msg的dataLen，放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack msg head error", err)
			break
		}
		// 根据dataLen，再次读取data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池即可
			c.MsgHandler.SendMsg2TaskQueue(req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Write Goroutine is running!]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	// 不断阻塞的等待channel的消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.GetTCPConnection().Write(data); err != nil {
				fmt.Println("Send data error, ", err, "conn Write exit！")
			}
		case <-c.ExitChan:
			// Reader退出，Writer也要退出
			return
		}
	}
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnId = ", c.ConnId)
	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 调用stop前的hook
	c.TcpServer.CallOnConnStop(c)
	// 关闭socket连接
	c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- struct{}{}
	// 将当前链接充ConnMgr中删除
	c.TcpServer.GetConnMgr().Remove(c.GetConnId())
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnId 获取当前连接模块的连接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 发送数据，将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包 |DataLen|Id|Data|
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgId = ", msgId)
		return errors.New("pack error msg")
	}
	// 将数据发送至Write（通过channel）
	c.msgChan <- binaryMsg
	return nil
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	// 添加一个属性
	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	// 删除属性
	delete(c.property, key)
}
