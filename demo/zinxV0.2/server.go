package main

import "zinx/znet"

func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.2]")
	// 2 启动server
	s.Serve()
}
