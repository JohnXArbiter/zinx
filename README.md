![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684332761525-dfcb88c8-5506-4b25-9ec7-aabe2bb62640.png#averageHue=%23ece9e5&clientId=u4201579d-43ee-4&from=paste&height=906&id=ud905d27f&originHeight=997&originWidth=1361&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=521465&status=done&style=none&taskId=u165a6a18-4819-4350-ba3e-b1ad8a033b2&title=&width=1237.2727004555636)

<a name="x3BJA"></a>
# 1 开始开发框架
使用接口抽象出我们zinx框架的功能
```go
package ziface

// IService 定义一个服务器接口
type IService interface {
	Start()
	Stop()
	Serve()
}
```
```go
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
	// TODO
}

func (s *Server) Stop() {
	// TODO
}
func (s *Server) Serve() {
	// TODO 
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

```

<a name="Qlc98"></a>
# 2 zinx V0.1 基础的server
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684411607350-cec1ad5c-2554-4c41-b9f0-bb56738b70a7.png#averageHue=%23d9d8d6&clientId=u1c8a1be3-2b6c-4&from=paste&height=263&id=ubef0aec8&originHeight=289&originWidth=817&originalType=binary&ratio=1&rotation=0&showTitle=false&size=117983&status=done&style=none&taskId=u228e25f1-1149-4713-beae-2c5678003d2&title=&width=742.727256629093)

一番操作后，我们又完成了一些我们的框架的server模块，先看看代码吧

Server() 是一个暴露给开发人员使用的函数，因为我们的服务肯定需要启动和停止，Server() 就可以负责，所以我们先写一个获取客户端连接然后将客户端发来的消息写回的功能
```go
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
```

1. 在start中，我们使用 `net.ResolveTCPAddr(network string, address string) (*TCPAddr, error)` 获取一个 TCP地址。
2. 然后我们需要对这个地址进行监听，调用`net.ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)`获取一个TCP监听器
3. 获取监听器以后，才能开始进行连接，使用`net.(l *TCPListener) AcceptTCP() (*TCPConn, error)`获取到连接，开始正式操作
<a name="tXTH3"></a>
# 2 zinx V0.2 简单的连接封装和业务绑定
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684413024596-0a1ad983-be31-4cfb-a5ac-e52bb3b0d647.png#averageHue=%23eeeeec&clientId=u1c8a1be3-2b6c-4&from=paste&height=452&id=u734a7e75&originHeight=497&originWidth=1302&originalType=binary&ratio=1&rotation=0&showTitle=false&size=253705&status=done&style=none&taskId=u187d18dc-55eb-4db8-b7a6-9d20625d642&title=&width=1183.636337981737)<br />我们在新增了项目目录的ziface包下新建icnnection.go文件，然后新建iconnection接口，是为了保存客户端和服务端长连接，我们需要将这个连接保存起来。<br />在项目目录下的znet包下去实现我们的接口，那Connection结构体为：
```go
// Connection 连接模块
type Connection struct {
	Conn      *net.TCPConn      // 当前连接的socket TCP套接字
	ConnID    uint32            // 连接的ID
	isClosed  bool              // 当前的连接状态
	handleAPI ziface.HandleFunc // 当前连接所绑定的处理业务方法API
	ExitChan  chan bool         // 告诉当前连接已经退出/停止 channel
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackApi,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}
```

每个连接都有自己的业务，会在Start()函数中开启，我们这里先写一个读业务的方法：
```go
// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}
		// 调用当前连接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID ", c.ConnID, " handle is error ", err)
			break
		}
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
```

那在zinx V0.2我们就完成对connection的编写，回到server，我们利用写好的connection更改一下server。我们将之前的写回操作替换成了去执行Connection中指定的HandlerFunc，只不过我们这也没有具体的HandlerFunc，我们就随便编写一个CalBackToClient
```go
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
		listener, err := net.ListenTCP(s.IPVersion, addr) // (*TCPListener, error)
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
				fmt.Println("Accept err", err)
				continue
			}
			// 客户端已经建立连接，开始做业务

			// 将处理新连接的业务方法 和 conn 进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// CallBackToClient 定义当前客户端链接的所绑定handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显的业务
	fmt.Println(" [Conn Handle] Cal1backToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
```

<a name="lpFCD"></a>
# 3 Zinx V0.3 基础router
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684494432051-769e74e5-b5e8-4287-a7ba-a126f83f2a7a.png#averageHue=%23e4e4e0&clientId=u99338fb1-3076-4&from=paste&height=797&id=u0f23a648&originHeight=877&originWidth=1579&originalType=binary&ratio=1&rotation=0&showTitle=false&size=728612&status=done&style=none&taskId=u1b1311e3-976c-4cb5-8510-ae196b41830&title=&width=1435.4545143419068)
<a name="KfJUc"></a>
## 3.1 封装请求Request
现在我们就给用户提供一个自定义的conn处理业务的接口吧，很显然，我们不能把业务处理业务的方法绑死在 type HandFunc func( *net.TCPConn，[ ]byte，int) error 这种格式中，我们需要定一些interface{}来让用户填写任意格式的连接处理业务方法。<br />那么，很显然func是满足不了我们需求的，我们需要再做几个抽象的接口类。

我们现在需要把客户端请求的连接信息和请求的数据，放在一个叫Request的请求类里，这样的好处是我们可以从Request里得到全部客户端的请求信息，也为我们之后拓展框架有一定的作用，一旦客户端有额外的含义的数据信息，都可以放在这个Request里。可以理解为每次客户端的全部请求数据，Zinx都会把它们一起放到一个Request结构体里。
```go
package ziface

/*
	IRequest接口：
	实际上是把客户端请求的连接信息和请求的数据包装到了一个Request中
*/
type IRequest interface {
	GetConnection() IConnection // 得到当前连接
	GetData() []byte            // 得到请求的消息数据
}
```
```go
package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好的连接
	data []byte             // 客户端请求的数据
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.data
}
```
<a name="OUjv6"></a>
## 3.2 抽象的IRouter
将我们封装的请求通过相应路由来处理，所以对应的路由中就需要三个handle
```go
package ziface

/*
	路由接口：
	路由里的数据都是IRequest
*/
type IRouter interface {
	PreHandle(request IRequest)  // 在处理conn业务之前的方法hook
	Handle(request IRequest)     // 在处理conn业务的主方法hook
	PostHandle(request IRequest) // 在处理conn业务之后的方法hook
}
```
因为handle肯定是使用人员去订制的，所有写一个BaseRouter而不是写死
```go
package znet

import "zinx/ziface"

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

/*	这里之所以BaseRouter的方法都为空
	是因为有的Router不希望有PreHandle、PostHandle这两个业务
	所以Router全部继承BaseRouter的好处就是，不需要实现PreHandle，PostHandle
*/
// PreHandle 在处理conn业务之前的方法hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {
}

// Handle 在处理conn业务的主方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {
}

// PostHandle 在处理conn业务之后的方法hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {
}
```
<a name="xtJp0"></a>
## 3.2 集成router
之前我们在connection里面调用当前连接的所绑定的业务，就一个写死的写回函数，现在我们改成上面写的IRequest接口，也就是封装的request。结构体和NewConnection这个函数，一样去掉之前的handlerFunc，改成新添加的IRouter接口
```go
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
		buf := make([]byte, 512)
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
```
server中也需要加入router：
```go
type Server struct {
	Name      string         // 服务器的名称
	IPVersion string         // 服务器绑定的ip版本
	IP        string         // 服务器监听的ip
	Port      int            // 服务器监听的端口
	Router    ziface.IRouter // 当前的Server添加一个router，server注册的连接对应的处理业务
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IService {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9000,
		Router:    nil,
	}
	return s
}
```
<a name="TyMJx"></a>
# 4 Zinx V0.4 全局配置
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684499229308-4d469560-e72e-4475-ac06-dcb71a5b800f.png#averageHue=%23ebece9&clientId=u99338fb1-3076-4&from=paste&height=256&id=u23781388&originHeight=282&originWidth=1690&originalType=binary&ratio=1&rotation=0&showTitle=false&size=308347&status=done&style=none&taskId=u8f237499-bede-4119-b423-3c4d5498052&title=&width=1536.363603063852)<br />随着架构逐步的变大，参数就会越来越多，为了省去我们后续大频率修改参数的麻烦，接下来Zinx需要做一个加载配置的模块，和一个全局获取Zinx参数的对象。

我们先做一个简单的加载配置模块，要加载的配置文件的文本格式，就选择比较通用的json格式，配置信息暂时如下：
```go
{
    "Name": "demo server",
    "Host": "127.0.0.1",
    "TcpPort": 7777,
	"MaxConn": 3
}
```

在utils下编写全局配置<br />定义GlobalObject结构体用于存储使用者编写的配置信息。用GlobalObject暴露给外面使用。init函数用来进行初始化：有默认值，然后去读取使用者的
```go
package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

var GlobalObject *GlobalObj

/*
	存储一切有关Zinx的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由使用人员进行配置
*/
type GlobalObj struct {
	TcpServer ziface.IService // 当前Zinx全局的Server对象
	Host      string          // 当前服务器主机监听的IP
	TcpPort   int             // 当前服务器主机监听的端口号
	Name      string          // 当前服务器的名称

	// Zinx
	Version        string // 当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 // 当前Zinx数据包的最大值
}

// 提供一个init函数，初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载，这是默认的值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4092,
	}

	// 应该尝试从conf/zinx.json去加载自定义的参数
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("config/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
```
将其他地方替换，比如这里获取server的时候
```go
// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IService {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
```

<a name="a62gu"></a>
# ⭐5 Zinx V0.5 消息模块与TCP数据包
<a name="a9iwq"></a>
## message
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684586091726-a38768eb-8b06-4e9e-bf5d-cfd4e255d109.png#averageHue=%23fafafa&clientId=u5a98cda1-5dbe-4&from=paste&height=521&id=u80b7b3a6&originHeight=573&originWidth=1315&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=275200&status=done&style=none&taskId=ud545da5a-c298-4819-b760-760adbb87d8&title=&width=1195.4545195437665)<br />接下来我们再对zinx做一个简单的升级，现在我们把服务器的全部数据都放在一个Request里，当前的Request结构如下∶
```go
type IRequest interface {
	GetConnection() IConnection // 得到当前连接
	GetData() []byte            // 得到请求的消息数据
}
```
很明显，现在是用一个[ ]byte来接受全部数据，又没有长度，又没有消息类型，这不科学。怎么办呢?我们现在就要自定义一种消息类型，把全部的消息都放在这种消息类型里。
```go
package ziface

/*
将请求的消息封装到一个Message中
*/
type IMessage interface {
	GetMsgId() uint32  // 获取消息的ID
	GetMsgLen() uint32 // 获取消息的长度
	GetData() []byte   // 获取消息的内容
	SetMsgId(uint32)   // 设置消息的ID
	SetData([]byte)    // 设置消息的内容
	SetDataLen(uint32) // 设置消息的长度
}
```
接口实现：<br />要解决我们的问题，所以我们的结构体肯定要有一定的属性：id、datalen、data。消息的长度下面我们解决TCP封包拆包的时候就会用到
```go
package znet

type Message struct {
	Id      uint32 // 消息的Id
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// GetMsgId 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
```
<a name="bY3iw"></a>
## TLV格式封包拆包
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684586386837-5c1427cf-4fc5-4c2f-beb3-a003b566b499.png#averageHue=%23ebebe8&clientId=u5a98cda1-5dbe-4&from=paste&height=185&id=u46ace52a&originHeight=204&originWidth=1297&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=172751&status=done&style=none&taskId=ufb563ee6-0f12-4099-82d8-8b279b9f76b&title=&width=1179.0908835348025)<br />![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684552685618-8d27be4c-91c0-4159-9105-06aca63e6486.png#averageHue=%23f3edec&clientId=u19ce8443-27b5-4&from=paste&height=693&id=u49cbec62&originHeight=693&originWidth=1415&originalType=binary&ratio=1&rotation=0&showTitle=false&size=295495&status=done&style=none&taskId=u371aeb1c-ac90-42a5-b04a-22135d58886&title=&width=1415)<br />由于Zinx也是TCP流的形式传播数据，难免会出现消息1和消息2一同发送，那么zinx就需要有能力区分两个消息的边界，所以zinx此时应该提供一个统一的拆包和封包的方法。在发包之前打包成如上图这种格式的有head和body的两部分的包，在收到数据的时候分两次进行读取，先读取固定长度的head部分，得到后续Data的长度，再根据DataLen读取之后的body。这样就能够解决粘包的问题了。
```go
package ziface

/*
封包、拆包模块
直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/
type IDataPack interface {
	GetHeadLen() uint32                // 获取包的头的长度方法
	Pack(msg IMessage) ([]byte, error) // 封包方法
	Unpack([]byte) (IMessage, error)   // 拆包方法
}
```
接口实现：<br />我们设定一个包的head长度是8字节，因为我们的message中定义的id为uint32是4位，datalen也是uint2 4位，所以Unpack我们就先读一个包的前8字节（只读包头），获取id和长度。然后根据这个长度再用一个buffer去读后面的数据
```go
package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// DataPack 封包、拆包的具体模块
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法
func (*DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节） + Id uint32（4字节）
	return 8
}

// Pack 封包方法
// |DataLen|MsgId|Data|
func (*DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 1 创建一个存放bytes字节的缓存
	dataBuffer := bytes.NewBuffer([]byte{})
	// 2 将dataLen写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 3 将MsgId写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 4 将data数据写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// Unpack 拆包方法（将包的Head信息读出来，之后再根据head的信息的data长度，再进行一次读）
func (*DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 1 创建一个从输入二进制数据的ioReader
	dataBuffer := bytes.NewReader(binaryData)
	// 2 只解压head信息，得到DataLen和MsgId
	msg := &Message{}
	// 3 读DataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 4 判断DataLen是否已经超出了最大允许包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}
	// 5 读MsgId
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	return msg, nil
}
```
然后我们编写一个测试文件还看看效果，也能加深一下理解：
```go
package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	// 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				// 处理客户端的请求
				// 拆包
				dp := NewDataPack()
				for {
					// 第一次读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					headMsg, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack head error")
						return
					}
					// 说明有数据，进行第二次读取
					if headMsg.GetDataLen() > 0 {
						// 第二次读，先把head的datalen再开始读data
						bodyData := make([]byte, headMsg.GetDataLen())
						_, err = io.ReadFull(conn, bodyData)
						if err != nil {
							fmt.Println("read body error")
							return
						}
						fmt.Println("----> Recv MsgId:", headMsg.GetMsgId(), " DataLen:", headMsg.GetDataLen(), " Data:", string(bodyData))
					}
				}
			}(conn)
		}
	}()

	// 客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}
	// 创建一个封包对象
	dp := NewDataPack()
	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'n', 'i', 'h', 'a', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
```
服务端：

- 就是监听相应的ip和端口，然后获得链接。接下来就使用我们的DataPack去拆包，可以看到，我们先定义了一个八字节长度的headData的字节数组，用它去读链接conn的前八个字节，得到的结果就是包头，然后用Unpack拆包，得到后面的数据长度。然后同样的，使用数据长度的bodyData在链接中去读这个长度的字节数，就拿到了数据

客户端：

- 连接相应的ip和端口，然后创建两个消息，指定长度、id、数据，然后都有Pack去封包，两个包进行粘包也就是放在一个字节数组中，然后发送

可以正常得到结果<br />![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684640265608-ea095ef9-9c49-4013-8222-3ca6b5af08b6.png#averageHue=%23222529&clientId=u61fb8cbd-ed96-4&from=paste&height=76&id=u901d18c2&originHeight=84&originWidth=445&originalType=binary&ratio=1&rotation=0&showTitle=false&size=10964&status=done&style=none&taskId=u0e2a246b-dd98-4929-967d-f62f0255b14&title=&width=404.54544577716814)
<a name="VUkjP"></a>
## 将消息封装到Zinx中
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1684640570263-645fa536-cbb4-4bbf-b49d-476cbd6d238c.png#averageHue=%23e7e9e7&clientId=u61fb8cbd-ed96-4&from=paste&height=191&id=ucbe750ca&originHeight=210&originWidth=1362&originalType=binary&ratio=1&rotation=0&showTitle=false&size=185540&status=done&style=none&taskId=u5f225c9a-ee18-4a51-8c9b-97b8e3de4ba&title=&width=1238.1817913449506)<br />之前我们的irequest的实现中只有一个connection和[]byte。那么在上面我们也已经写好了message，是时候替换了<br />irequest的实现：
```go
package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好的连接
	msg  ziface.IMessage    // 客户端请求的数据
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 消息的id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
```
我们是通过connection去调用request，request现在又封装了message，所以connection里面也要做更改<br />将之前读取客户端发来的数据包进行现在拆包再获取数据，并增加了SendMsg函数封装了发送数据包的过程
```go
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
			_, err := io.ReadFull(c.GetTCPConnection(), data)
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
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
		// 从路由中，找到注册绑定的Conn对应的router调用
	}
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
	// 将数据发送至客户端
	if _, err = c.GetTCPConnection().Write(binaryMsg); err != nil {
		fmt.Println("Write error msgId = ", msgId, " error: ", err)
		return errors.New("conn write error")
	}
	return nil
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
```
<a name="jrA9B"></a>
# 6 Zinx V0.6 多路由模式
<a name="G9j2D"></a>
## 6.1 msgHandler管理路由handler
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685453697045-cbb3e3ba-0319-40dd-a5b1-c0bd7a2f66d8.png#averageHue=%23dedfdb&clientId=ue802ea0c-895a-4&from=paste&height=167&id=ua027c824&originHeight=184&originWidth=1127&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=129820&status=done&style=none&taskId=u50e32a7c-d5b1-4996-b2ce-9dd05e95f5e&title=&width=1024.5454323390304)<br />之前我们没有做多路由多handler的功能，我们是直接在server中添加了irouter成员，connection中也添加了irouter成员<br />那么下面我们在connection通过路由调用相应的handler就只是这样
```go
// 得到当前conn数据的Request请求数据
req := &Request{
        conn: c,
        msg:  msg,
	}
// 执行注册的路由方法
go func(request ziface.IRequest) {
        c.Router.PreHandle(request)
        c.Router.Handle(request)
        c.Router.PostHandle(request)
    }(req)
```
现在，我们新封装一个多handler的模块<br />用map去存储相应msgId对应的handler
```go
package ziface

/*
消息管理抽象
*/
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 调度/执行对应的Router消息处理方法
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
}
```
接口实现：
```go
package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgId对应的处理方法
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgId
	msgId := request.GetMsgId()
	handler, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " NOT FOUND! Need Register First!")
	}
	// 2 根据msgId，调度相应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// msgId已经注册
		panic("repeat api, msgId: " + strconv.Itoa(int(msgId)))
	}
	// 添加绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router MsgId: ", msgId, " success")
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
```
<a name="LEvaI"></a>
## 6.2 集成到之前的代码中
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685455803542-7948522e-b771-4a01-b5cc-923d753c1aee.png#averageHue=%23e4e5e2&clientId=ue802ea0c-895a-4&from=paste&height=266&id=u15274c09&originHeight=293&originWidth=1335&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=276640&status=done&style=none&taskId=ua246ae05-5c90-4579-a563-bc18e8029d2&title=&width=1213.6363373315044)<br />先更改server，将之前irouter的单个路由替换成IMsgHandler这个可以存储多个handler，那么初始化和添加路由也需要调用这个接口中的函数
```go
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
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackeetSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
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
```
connection也是一样的，将之前的irouter替换掉，因为我们要在相应的connection中调用当前消息的handler
```go
// Connection 连接模块
type Connection struct {
	Conn       *net.TCPConn       // 当前连接的socket TCP套接字
	ConnID     uint32             // 连接的ID
	isClosed   bool               // 当前的连接状态
	ExitChan   chan bool          // 告诉当前连接已经退出/停止 channel
	MsgHandler ziface.IMsgHandler // 消息的管理MsgID和对应的处理业务API关系
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
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
			_, err := io.ReadFull(c.GetTCPConnection(), data)
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
```

<a name="m20BW"></a>
# 7 Zinx V0.7 消息读写分离![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685507383950-dc8009a8-d31f-43b1-89fc-707ebd17f3b0.png#averageHue=%23e2e3df&clientId=u4234adb3-e16c-4&from=paste&height=160&id=u01ded56b&originHeight=176&originWidth=700&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=82135&status=done&style=none&taskId=ud57180a1-f121-47bb-8ac4-00e5311e5d2&title=&width=636.3636225708264)
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685455534856-ae23ea38-0036-4433-9a1e-8f43ed92f7b8.png#averageHue=%23fbfbfb&clientId=ue802ea0c-895a-4&from=paste&height=933&id=u62903c4f&originHeight=1026&originWidth=1258&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=277384&status=done&style=none&taskId=uf6da841a-eccd-4681-86c7-db5eb9b1c2f&title=&width=1143.6363388487136)<br />上面一节，我们还能看到，在connection的reader中，我们通过接收客户端发送过来的数据后，直接给到了msgHandler中去处理，毕竟处理业务一般也会有响应的业务，我们就尝试将其分离<br />看到上面可爱的goroutine代表的就是我们会有至少三个协程去处理，一个server的主协程，然后就是reader和writer这两个协程分别去做读写，那写肯定是需要在读完成之后的，所以我们通过channel去实现协程通信

那么在代码中，我们在Connection结构体中，用上了ExitChan，并加入了msgChan，一个用来reader去通知writer结束，还有一个就是用来传递数据。<br />在一开始调用connection时的Start函数，我们就直接开启这两个协程，writer阻塞的不断去取数据，reader也是不断的去接收数据，然后在msgHandler中我们会去调用SendMsg那时候就会向channel给数据，writer就可以写了
```go
package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
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
	fmt.Println("[Reader Goroutine is running!]")
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
```
<a name="KOm0D"></a>
# ⭐8 Zinx V0.8 消息队列和多任务处理机制
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685505633397-b2c31832-efcc-4e2d-978c-05f89b9cc3be.png#averageHue=%23f8f3f2&clientId=u4234adb3-e16c-4&from=paste&height=718&id=u542a65d0&originHeight=790&originWidth=1358&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=323476&status=done&style=none&taskId=u6cb40c2c-7a4b-40a2-acb6-c82ee835204&title=&width=1234.545427787403)<br />根据0.7我们加入的读写两个goroutine，会出现一个问题：<br />我们在reader中，去执行了相应的msgHandler，这也是一个goroutine。所以当并发量上来的时候，比如说一下来了十万个不同链接`（注意我们就是根据链接，所以说一个）`，那我们就会创建30W的goroutine（10W个reader，10W个writer，10W个handler），由于我们还是使用channel进行reader和writer通信的（也就是进入到handler后，调用sendMsg发送数据到channel，才能让writer继续），那在msgHandler执行完之前，writer是阻塞的，不占用cpu，msgHandler也执行得很快，但是一个goroutine执行完后切换到另一个goroutine的线程切换开销还是存在，而且在并发量如此大的情况下

现在我们需要做的是：
<a name="QpVjv"></a>
## 8.1 创建一个消息队列
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685519940875-19e9f5c1-6630-421e-9f94-b83ef03dada4.png#averageHue=%23e2e3e0&clientId=u4234adb3-e16c-4&from=paste&height=67&id=u44affead&originHeight=74&originWidth=1124&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=69195&status=done&style=none&taskId=ub60b785c-f4fe-448c-a091-8aa79a10835&title=&width=1021.8181596708697)<br />在MsgHandler中添加TaskQueue，一个IRequest类型的channel切片，切片的长度就是工作池的大小，切片里的channel才是真正的消息队列
```go
package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgId对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作Worker池的worker数量
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgId
	msgId := request.GetMsgId()
	handler, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " NOT FOUND! Need Register First!")
	}
	// 2 根据msgId，调度相应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// msgId已经注册
		panic("repeat api, msgId: " + strconv.Itoa(int(msgId)))
	}
	// 添加绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router MsgId: ", msgId, " success")
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}
```
我们在配置这里也加入工作池大小的设定
```go
type GlobalObj struct {
	TcpServer ziface.IService // 当前Zinx全局的Server对象
	Host      string          // 当前服务器主机监听的IP
	TcpPort   int             // 当前服务器主机监听的端口号
	Name      string          // 当前服务器的名称

	// Zinx
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 // 当前Zinx数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker池的goroutine数量
	MaxWorkerTaskLen uint32 // 允许用户最多开辟多少个Worker
}

// 提供一个init函数，初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载，这是默认的值
	GlobalObject = &GlobalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		MaxConn:          1000,
		MaxPackageSize:   4092,
		WorkerPoolSize:   10,   // worker工作池的队列的个数
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的任务的数量的最大值
	}

	// 应该尝试从conf/zinx.json去加载自定义的参数
	//GlobalObject.Reload()
}
```
<a name="bly89"></a>
## 8.2 创建多任务worker的工作池并启动
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685519922447-b9d220ad-f5a4-47a1-b2bc-3e6afe2c3623.png#averageHue=%23e2e3e0&clientId=u4234adb3-e16c-4&from=paste&height=132&id=u118c2ed9&originHeight=145&originWidth=1209&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=123565&status=done&style=none&taskId=u68229b49-99ab-4a3c-a8cb-629dd74f6df&title=&width=1099.0908852687558)<br />调用StartWorkerPool函数，就是开启我们的多消息队列去处理消息的模式（只需调用一次即可），他就会去调用StartOneWorker函数，就是开启设定的数量个数个channel，然后不断去监听过来的request，有request到来，就去执行相应的msgHandler
```go
package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgId对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作Worker池的worker数量
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgId
	msgId := request.GetMsgId()
	handler, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " NOT FOUND! Need Register First!")
	}
	// 2 根据msgId，调度相应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// msgId已经注册
		panic("repeat api, msgId: " + strconv.Itoa(int(msgId)))
	}
	// 添加绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router MsgId: ", msgId, " success")
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// StartWorkerPool 启动一个Worker工作池（开启工作池的动作只发生一次，一个zinx只能有一个worker工作池）
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 为当前的worker对应的channel消息队列初始化，第i个worker就用第i个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker，阻塞等待消息从channel传递过来
		go mh.StartOneWork(i)
	}
}

// StartOneWork 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWork(workerId int) {
	fmt.Println("Worker Id = ", workerId, " is started ...")
	// 不断地阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-mh.TaskQueue[workerId]:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsg2TaskQueue(request ziface.IRequest) {
	// 1 将消息平均分给不通过的worker
	// 1.1 根据客户端的链接id来进行分配（同一个链接的request全都会在这个队列里）
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnID(),
		" request MsgId = ", request.GetMsgId(),
		" to WorkerId = ", workerId)
	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerId] <- request
}
```
<a name="e0L0g"></a>
## 8.3 将之前的发送消息，全部改成 把消息发送给 消息队列和worker工作池来处理
那上面那一步的request怎么来的呢，之前我们是connection中的reader中直接调用的msgHandler去处理的，现在我们还是在connection的reader中获得request，但是现在我们发到我们的msgHandler.SendMsg2TaskQueue函数通过channel给到刚刚开启但阻塞的StartOneWorker函数再去处理<br />connection.go：
```go
// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running!]")
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池即可
			c.MsgHandler.SendMsg2TaskQueue(req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}
```
刚刚也是说了，开启workerPool也只需要一次，我们就在server中去开启就行了
```go
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
```

<a name="XmRCd"></a>
# 9 Zinx V0.9 连接管理
<a name="KuLUp"></a>
## 9.1 定义链接管理模块
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685618232841-e49384ab-2d2d-4664-b9a7-4aa82d0464a0.png#averageHue=%23eeeeeb&clientId=u7b515251-745d-4&from=paste&height=309&id=u2c10ce54&originHeight=340&originWidth=1035&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=183959&status=done&style=none&taskId=ude02919e-c5ca-457c-9a2d-d3e6ccee048&title=&width=940.909070515436)<br />现在我们要为Zinx框架增加链接个数的限定，如果超过一定量的客户端个数，Zinx为了保证后端的及时响应，而拒绝链接请求。<br />我们需要一个map[connId]iconnection，来通过链接的id保存链接，在iconnmanager这个接口中定义了一系列的操作<br />定义接口：
```go
package ziface

/*
* 连接管理模块抽象层
 */

type IConnManager interface {
	Add(conn IConnection)                   // 添加链接
	Remove(connId uint32)                   // 删除链接
	Get(connId uint32) (IConnection, error) // 根据connId获取链接
	Len() int                               // 得到当前链接总数
	ClearConn()                             // 清除并中止所有的连接
}
```
接口实现：<br />因为map不是线程安全的，所以我们在manager里面加入了读写锁，保证读写并发安全。<br />比如在Add和Remove中，就需要先加上写锁，defer释放
```go
package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// ConnManager 链接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的链接集合
	connLock    sync.RWMutex                  // 保护链接集合的读写锁
}

// Add 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 将conn加入到ConnManager中
	cm.connections[conn.GetConnId()] = conn

	fmt.Println("connId = ", conn.GetConnId(), " add to ConnManager successfully: conn num = ", cm.Len())
}

// Remove 删除链接
func (cm *ConnManager) Remove(connId uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, connId)
	fmt.Println("connId = ", connId, " remove from ConnManager successfully: conn num = ", cm.Len())
}

// Get 根据connId获取链接
func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		// 找到
		return conn, nil
	}
	return nil, errors.New("connection NOT FOUND")
}

// Len 得到当前链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除并中止所有的连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除conn并停止conn的工作
	for connId, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connId)
	}
	fmt.Println("Clear All connections success! conn num = ", cm.Len())
}
```
<a name="BkNAG"></a>
## 9.2 与server、connection模块关联
我们先更改server模块<br />![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685618428600-2e72ee12-ac5f-4650-9a79-8327d4b63327.png#averageHue=%23e1e2df&clientId=u7b515251-745d-4&from=paste&height=117&id=u4f48e2ba&originHeight=129&originWidth=811&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=83930&status=done&style=none&taskId=u07e53c19-6df7-4db9-964e-52cb2aa3c4c&title=&width=737.2727112927716)<br />在接口中添加一个GetConnMgr的函数，也就是说，一个server对应一个他自己的connMgr，而每个connMgr中就会有很多个connection
```go
// IServer 定义一个服务器接口
type IServer interface {
	Start()                                 // 启动服务器
	Stop()                                  // 停止服务器
	Serve()                                 // 运行服务器
	AddRouter(msgId uint32, router IRouter) // 路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	GetConnMgr() IConnManager               // 获取当前server的连接管理器
}
```
在server实现中加入connMgr，然后在链接过来时，我们会做一个判断，判断当前已有的链接数有没有大于等于最大链接数
```go
type Server struct {
	Name       string              // 服务器的名称
	IPVersion  string              // 服务器绑定的ip版本
	IP         string              // 服务器监听的ip
	Port       int                 // 服务器监听的端口
	MsgHandler ziface.IMsgHandler  // 当前server的消息管理模块，用来绑定和对应的处理业务API关系
	ConnMgr    ziface.IConnManager // 该server的链接管理器
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

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnMgr(),
	}
}
```
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685618814706-b2ce813e-a84b-49ee-9678-ffe302eed9b0.png#averageHue=%23e1e4e2&clientId=u7b515251-745d-4&from=paste&height=104&id=u7b916191&originHeight=114&originWidth=953&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=120086&status=done&style=none&taskId=u140ac0cc-083f-440c-a55f-5a22b6f4b96&title=&width=866.3636175857107)<br />对应connection，我们添加了TcpServer这个成员，也就是server，因为每个server对应着多个connection，那我们也应该从connection能找到server，所以在初始化时，我们将TcpServer绑定好，并且将自己（connection）加入server的connMgr中
```go
// Connection 连接模块
type Connection struct {
	TcpServer  ziface.IServer     // 当前connection属于哪个server
	Conn       *net.TCPConn       // 当前连接的socket TCP套接字
	ConnId     uint32             // 连接的ID
	isClosed   bool               // 当前的连接状态
	ExitChan   chan struct{}      // Reader告诉Writer去退出 channel
	msgChan    chan []byte        // 无缓冲管道，用于读、写goroutine之间的消息通信
	MsgHandler ziface.IMsgHandler // 消息的管理MsgID和对应的处理业务API关系
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
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}
```
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685622048655-8d1a4504-3dd4-4cff-9dfa-3f7128de3015.png#averageHue=%23e1e1de&clientId=u7b515251-745d-4&from=paste&height=376&id=u9f1d1372&originHeight=414&originWidth=901&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=264211&status=done&style=none&taskId=ufc01f887-3315-4aa5-9ebe-55da2c66123&title=&width=819.0908913375921)
```go
package ziface

// IServer 定义一个服务器接口
type IServer interface {
	Start()                                      // 启动服务器
	Stop()                                       // 停止服务器
	Serve()                                      // 运行服务器
	AddRouter(msgId uint32, router IRouter)      // 路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	GetConnMgr() IConnManager                    // 获取当前server的连接管理器
	SetOnConnStart(func(connection IConnection)) // 注册OnConnStart hook的方法
	SetOnConnStop(func(connection IConnection))  // 注册OnConnStop hook的方法
	CallOnConnStart(connection IConnection)      // 调用OnConnStart hook的方法
	CallOnConnStop(connection IConnection)       // 调用OnConnStop hook的方法
}
```
```go
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
```

<a name="MMG2z"></a>
# 10 Zinx V1.0 添加链接属性
![image.png](https://cdn.nlark.com/yuque/0/2023/png/34740401/1685624172616-f4f4c4e5-728b-4e17-9059-928161928e28.png#averageHue=%23ecedea&clientId=u7b515251-745d-4&from=paste&height=169&id=u81a8a185&originHeight=186&originWidth=862&originalType=binary&ratio=1.100000023841858&rotation=0&showTitle=false&size=114941&status=done&style=none&taskId=uac9c92da-4c04-4c1d-9d78-3f488f9ec74&title=&width=783.6363466515032)<br />用一个map保存自定义的属性，不用多说了吧
```go
// IConnection 定义连接模块的抽象层
type IConnection interface {
	Start()                                      // 启动连接，让当前的连接准备开始工作
	Stop()                                       // 停止连接，结束当前连接的工作
	GetTCPConnection() *net.TCPConn              // 获取当前连接绑定的socket conn
	GetConnId() uint32                           // 获取当前连接模块的连接ID
	RemoteAddr() net.Addr                        // 获取远程客户端的 TCP状态 IP port
	SendMsg(msgId uint32, data []byte) error     // 发送数据，将数据发送给远程的客户端
	SetProperty(key string, value interface{})   // 设置链接属性
	GetProperty(key string) (interface{}, error) // 获取链接属性
	RemoveProperty(key string)                   // 移除链接属性
}

```
```go
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
```
