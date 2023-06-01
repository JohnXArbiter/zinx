package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// PreHandle 在处理conn业务之前的方法hook
func (br *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle 在处理conn业务的主方法hook
func (br *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping..."))
	if err != nil {
		fmt.Println("call back ping ping ping... error")
	}
}

// PostHandle 在处理conn业务之后的方法hook
func (br *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Cal1 Router AfterHandle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer()
	// 2 添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3 启动server
	s.Serve()
}
