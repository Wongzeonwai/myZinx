package ziface

type IMessage interface {
	GetMsgID() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgID(id uint32)
	SetMsgLen(l uint32)
	SetData(data []byte)
}
