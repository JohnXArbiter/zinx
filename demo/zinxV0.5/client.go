package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	// 1 直接连接远程服务器，得到conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		// 发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.5 client test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error, ", err)
			return
		}
		// 服务器应该回复一个message数据， MsgId：1 ping ping ping...

		// 先读取流中head部分，得到id和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error")
			break
		}
		// 再根据dataLen进行第二次读取，将data读出来
		msg, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack head error")
			break
		}
		// 说明有数据，进行第二次读取
		if msg.GetDataLen() > 0 {
			// 第二次读，先把head的datalen再开始读data
			msg.SetData(make([]byte, msg.GetDataLen()))
			_, err = io.ReadFull(conn, msg.GetData())
			if err != nil {
				fmt.Println("read body error")
				return
			}
			fmt.Println("----> Recv MsgId:", msg.GetMsgId(), " DataLen:", msg.GetDataLen(), " Data:", string(msg.GetData()))

		}
		// cpu阻塞
		time.Sleep(time.Second)
	}
}
