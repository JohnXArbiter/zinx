package znet

type Message struct {
	Id      uint32 // 消息的Id
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// GetMsgId 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetDataLen 获取消息的长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
