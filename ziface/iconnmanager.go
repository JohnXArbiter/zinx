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
