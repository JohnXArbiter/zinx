package ziface

/*
消息管理抽象
*/
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 调度/执行对应的Router消息处理方法
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
}
