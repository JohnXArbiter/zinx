package ziface

// IService 定义一个服务器接口
type IService interface {
	Start()
	Stop()
	Serve()
}
