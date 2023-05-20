package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// DataPack 封包、拆包的具体模块
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法
func GetHeadLen() uint32 {
	// DataLen uint32（4字节） + Id uint32（4字节）
	return 8
}

// Pack 封包方法
// |DataLen|MsgId|Data|
func Pack(msg ziface.IMessage) ([]byte, error) {
	// 1 创建一个存放bytes字节的缓存
	dataBuffer := bytes.NewBuffer([]byte{})
	// 2 将dataLen写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 3 将MsgId写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	// 4 将data数据写进dataBuffer中
	return dataBuffer.Bytes(), nil
}

// Unpack 拆包方法（将包的Head信息读出来，之后再根据head的信息的data长度，再进行一次读）
func Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 1 创建一个从输入二进制数据的ioReader
	dataBuffer := bytes.NewReader(binaryData)

	// 2 只解压head信息，得到DataLen和MsgId
	msg := &Message{}
	// 3 读DataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 4 判断DataLen是否已经超出了最大允许包长度
	if utils.GlobalObj.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObj.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}
	// 5 读MsgId
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	return msg, nil
}
