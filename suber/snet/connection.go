package snet

import (
	"errors"
	"fmt"
	"io"
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
		/*
		//读取客户端数据到buf
		buf := make([]byte,utils.GlobalConfigNow.MaxPackageSize)
		_,err := c.Conn.Read(buf)
		if err != nil{
			fmt.Println("c.Conn.Read(buf) err = ",err)
			continue
		}
		*/

		//创建一个拆包解包的对象
		dp:=NewDataPack()

		//读取客户端的Msg Head 二进制流
		headData := make([]byte,dp.GetHeadLen())
		if _,err:=io.ReadFull(c.GetTCPConnection(),headData);err!=nil{
			fmt.Println("o.ReadFull(c.GetTCPConnection() err",err)
			break
		}

		//拆包，等到MsgID和MsgLen,放到Msg消息中
		msg,err:=dp.UnPack(headData)
		if err != nil{
			fmt.Println("dp.UnPack(headData) err = ",err)
			break
		}
		var data []byte
		if msg.GetMsgLen()>0{
			data = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(c.GetTCPConnection(),data);err!=nil{
				fmt.Println("io.ReadFull(c.GetTCPConnection(),data) err = ",err)
				break
			}
		}
		msg.SetMsgData(data)


		//得到dataLen,再次读取Data,放到msg Data中


		//得到当前Conn数据的request请求
		req := &Request{
			conn:c,
			msg:msg,
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

//提供一个SendMsg方法，将我们要发给客户端的数据先进行封包，再发送
func (c *Connection)SendMsg(msgId uint32,data []byte)error{
	if c.IsCLose == true{
		return errors.New("connection closed when send msg")
	}

	//将data进行封包MsgLen| MsgID | data

	dp:=NewDataPack()

	binMsg,err:=dp.Pack(NewMsgPackage(msgId,data))
	if err!=nil {
		fmt.Println("dp.Pack(NewMsgPackage(msgId,data)) err = ",err)
		return errors.New("pack error msg")
	}

	//将数据发送给客户端
	if _,err:=c.Conn.Write(binMsg);err != nil{
		fmt.Println("c.Conn.Write(binMsg) err = ",err)
		return errors.New("c.Conn.Write err")
	}
	return nil
}
