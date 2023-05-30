package ziface

/*
将请求的消息封装到一个Message中
*/
type IMessage interface {
	GetMsgId() uint32   // 获取消息的ID
	GetDataLen() uint32 // 获取消息的长度
	GetData() []byte    // 获取消息的内容
	SetMsgId(uint32)    // 设置消息的ID
	SetData([]byte)     // 设置消息的内容
	SetDataLen(uint32)  // 设置消息的长度
}
