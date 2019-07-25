package siface

import "net"

//定义链接模块抽象层

type IConnection interface {
	//启动链接 让当前的链接准备开始工作
	Start()
	//停止链接 结束当前链接的工作
	Stop()
	//获取当前链接所绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接的模块的ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr
	//发送数据 将数据发送到运程客户
	SendMsg(msgId uint32,data []byte) error
}

type HandleFunc func(*net.TCPConn,[]byte,int) error