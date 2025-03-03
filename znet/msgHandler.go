package znet

import (
	"fmt"
	"go-zinx/ziface"
	"strconv"
)

type MsgHandle struct {
	APIs map[uint32]ziface.IRouter
}

// 初始化、创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs: make(map[uint32]ziface.IRouter),
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
