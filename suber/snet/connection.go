package snet

import (
	"fmt"
	"net"
	"suber/siface"
)

/*
	链接模块
*/
type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//当前链接ID
	ConnID uint32

	//当前链接状态
	IsCLose bool

	//告知当前链接已经退出的/停止channel
	ExitChan chan bool

	//该链接处理的方法Router

	Router siface.IRouter

}

//初始化链接模块的方法

func NewConnection(conn *net.TCPConn,connID uint32,router siface.IRouter) *Connection  {
	c:=&Connection{
		Conn:conn,
		ConnID:connID,
		Router:router,
		IsCLose:false,
		ExitChan:make(chan bool,1),
	}
	return c
}

func (c *Connection)StartReader()  {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ",c.ConnID,"Reader is exit ,remote addr is",c.RemoteAddr().String())
	defer c.Stop()

	for{
		//读取客户端数据到buf 最大512byte
		buf := make([]byte,512)
		_,err := c.Conn.Read(buf)
		if err != nil{
			fmt.Println("c.Conn.Read(buf) err = ",err)
			continue
		}
		//得到当前Conn数据的request请求
		req := &Request{
			conn:c,
			data:buf,
		}
		go func(request siface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
		//调用路由，从路由中找到Conn对应的Router调用

	}
}
//启动链接 让当前的链接准备开始工作
func (c *Connection)Start(){
	fmt.Println("Connetion Start()..ConnID:",c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//TODO 从当前链接的写数据业务

}
//停止链接 结束当前链接的工作
func (c *Connection) Stop(){
	fmt.Println("Connetion Stop()....ConnID:",c.ConnID)
	if c.IsCLose == true{
		return
	}
	c.IsCLose = true
	//	回收资源
	c.Conn.Close()
	close(c.ExitChan)
}
//获取当前链接所绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn{
	return c.Conn
}
//获取当前链接的模块的ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}
//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}
//发送数据 将数据发送到运程客户端
func (c *Connection) Send(data []byte) error{
	return nil
}

