package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgId对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作Worker池的worker数量
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
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// StartWorkerPool 启动一个Worker工作池（开启工作池的动作只发生一次，一个zinx只能有一个worker工作池）
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 为当前的worker对应的channel消息队列初始化，第i个worker就用第i个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker，阻塞等待消息从channel传递过来
		go mh.StartOneWork(i)
	}
}

// StartOneWork 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWork(workerId int) {
	fmt.Println("Worker Id = ", workerId, " is started ...")
	// 不断地阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-mh.TaskQueue[workerId]:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsg2TaskQueue(request ziface.IRequest) {
	// 1 将消息平均分给不通过的worker
	// 1.1 根据客户端的链接id来进行分配（同一个链接的request全都会在这个队列里）
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnID(),
		" request MsgId = ", request.GetMsgId(),
		" to WorkerId = ", workerId)
	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerId] <- request
}
