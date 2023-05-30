package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个MsgId对应的处理方法
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgId
	msgId := request.GetMsgId()
	handler, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " NOT FOUND! Need Register First!")
	}
	// 2 根据msgId，调度相应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// msgId已经注册
		panic("repeat api, msgId: " + strconv.Itoa(int(msgId)))
	}
	// 添加绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router MsgId: ", msgId, " success")
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
