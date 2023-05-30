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
