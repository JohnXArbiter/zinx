package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	Conn       *net.TCPConn       // 当前连接的socket TCP套接字
	ConnID     uint32             // 连接的ID
	isClosed   bool               // 当前的连接状态
	ExitChan   chan struct{}      // Reader告诉Writer去退出 channel
	msgChan    chan []byte        // 无缓冲管道，用于读、写goroutine之间的消息通信
	MsgHandler ziface.IMsgHandler // 消息的管理MsgID和对应的处理业务API关系
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan struct{}, 1),
	}
}

// Start 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// 启动从当前连接的读数据业务
	go c.StartReader()
	// 启动从当前连接的写数据业务
	go c.StartWriter()
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
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
		// 从路由中，找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(req)
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
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- struct{}{}
	close(c.ExitChan)
	close(c.msgChan)
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
