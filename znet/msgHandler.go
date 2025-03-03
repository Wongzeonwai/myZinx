package znet

import (
	"fmt"
	"go-zinx/utils"
	"go-zinx/ziface"
	"strconv"
)

type MsgHandle struct {
	APIs map[uint32]ziface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
}

// 初始化、创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("no such handler")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.APIs[msgID]; ok {
		panic("duplicated msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.APIs[msgID] = router
	fmt.Println("add msgID = " + strconv.Itoa(int(msgID)))
}

// 启动一个worker工作池（只能发生一次，一个zinx框架只能有一个工作池）
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 1、当前worker对应的channel开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2、启动当前worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker id = ", workerID)

	// 不断阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息进来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给taskqueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1、将消息平均分配给不同的worker
	// 根据客户端建立的connID分配
	workID := request.GetConnection().GetID() % mh.WorkerPoolSize
	fmt.Println("send msg to worker id = ", workID)
	// 2、将消息发送给对应的taskqueue
	mh.TaskQueue[workID] <- request
}
