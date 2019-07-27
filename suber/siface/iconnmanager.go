package siface

/*
	链接管理模块
*/
type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//得到链接总数
	Len()int
	//得到某个连接
	Get(connID uint32)(IConnection,error)
	//清除并终止所有链接
	ClearConn()

}
