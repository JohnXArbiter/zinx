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
	//GlobalObject.Reload()
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
