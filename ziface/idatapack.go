package ziface

// 解决TCP粘包问题的封包拆包模块
type IDataPack interface {
	// 获取包长度方法
	GetHeadLen() uint32
	// 封包方法
	Pack(msg IMessage) ([]byte, error)
	// 拆包方法
	Unpack([]byte) (IMessage, error)
}
