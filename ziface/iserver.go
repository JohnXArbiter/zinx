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
