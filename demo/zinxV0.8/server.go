package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// Handle 在处理conn业务的主方法hook
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	// 先读取客户端的数据，在写回ping ping ping...
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping ping ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle 在处理conn业务的主方法hook
func (br *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle..")
	// 先读取客户端的数据，在写回ping ping ping...
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello welcome to Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer()
	// 2 添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 3 启动server
	s.Serve()
}
