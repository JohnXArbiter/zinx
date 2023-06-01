package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/utils"
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

func NewConnMgr() ziface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection, utils.GlobalObject.MaxConn),
	}
}
