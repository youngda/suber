package snet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"suber/siface"
	"suber/utils"
	"sync"
)

/*
	链接模块
*/
type Connection struct {
	// 当前conn属于哪个server
	Tcpserver siface.IServer
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//当前链接ID
	ConnID uint32

	//当前链接状态
	IsCLose bool

	//告知当前链接已经退出的/停止channel
	ExitChan chan bool

	// 无缓冲的管道，用于读写Gouroutine之间的消息通信
	msgChan chan []byte
	//消息的管理MsgID，和对应的处理业务API关系
	MsgHander siface.IMsgHandler

	//链接属性集合
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

//初始化链接模块的方法

func NewConnection(server siface.IServer,conn *net.TCPConn, connID uint32, msgHandler siface.IMsgHandler) *Connection {
	c := &Connection{
		Tcpserver:server,
		Conn:      conn,
		ConnID:    connID,
		MsgHander: msgHandler,
		IsCLose:   false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property: make(map[string]interface{}),
	}
	c.Tcpserver.GetConnMgr().Add(c)
	//将conn加入到connManager中
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit] ,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
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
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("o.ReadFull(c.GetTCPConnection() err", err)
			break
		}

		//拆包，等到MsgID和MsgLen,放到Msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("dp.UnPack(headData) err = ", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("io.ReadFull(c.GetTCPConnection(),data) err = ", err)
				break
			}
		}
		msg.SetMsgData(data)

		//得到dataLen,再次读取Data,放到msg Data中

		//得到当前Conn数据的request请求
		req := &Request{
			conn: c,
			msg:msg,
		}
		if utils.GlobalConfigNow.WorkerPoolSize>0{
			//已经开启了工作池机制，将消息发送给Worker工作池处理
			c.MsgHander.SendMsgToTaskQueue(req)
		}else{
			//调用路由，从路由中找到Conn对应的Router调用
			go c.MsgHander.DoMsgHandler(req)
		}

	}
}

/*
	写消息的Goroutine,专门发送给客户端消息的模块
*/
func (c *Connection) StartWriter() {
	//
	fmt.Println("[Writer Coroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer closed]")
	//阻塞等待channel的消息，进行回写客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error c.Conn.Write(data) err = ",err)
				return
			}
		case <-c.ExitChan:
			fmt.Println("writer c.ExitChan exit")
			return
		}
	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connetion Start()..ConnID:", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//从当前链接的写数据业务
	go c.StartWriter()

	//按照开发者传递来的 创建链接之后需要调用的处理业务，执行对应的hook
	c.Tcpserver.CallOnConnStart(c)
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Connetion Stop()....ConnID:", c.ConnID)
	if c.IsCLose == true {
		return
	}
	c.IsCLose = true

	//将当前链接从connMgr删除
	c.Tcpserver.GetConnMgr().Remove(c)



	//调用开发者传递来的 销毁链接之前需要调用的处理业务，执行对应的hook
	c.Tcpserver.CallOnConnStop(c)


	//	回收资源
	c.Conn.Close()
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前链接所绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接的模块的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供一个SendMsg方法，将我们要发给客户端的数据先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsCLose == true {
		return errors.New("connection closed when send msg")
	}

	//将data进行封包MsgLen| MsgID | data

	dp := NewDataPack()

	binMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("dp.Pack(NewMsgPackage(msgId,data)) err = ", err)
		return errors.New("pack error msg")
	}

	//将数据发送给chan
	c.msgChan <-binMsg

	return nil
}



//设置链接属性
func (c *Connection) SetProperty(key string,value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个链接属性

	c.property[key] = value
}
//获取链接属性
func (c *Connection) GetProperty(key string)(interface{},error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	//读取属性
	if value,ok:=c.property[key];ok{
		return value,nil
	}else {
		return nil,errors.New("no property found")
	}
}
//移除链接属性
func (c *Connection) RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//删除属性
	delete(c.property,key)

}

