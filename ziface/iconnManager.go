package ziface

type IConnManager interface {
	// 添加连接
	AddConn(conn IConnection)
	// 删除连接
	RemoveConn(conn IConnection)
	// 根据connID获取连接
	GetConnCount(connID uint32) (IConnection, error)
	// 得到当前连接总数
	GetConnLen() int
	// 清除并终止所有连接
	ClearConn()
}
