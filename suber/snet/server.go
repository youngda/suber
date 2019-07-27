package snet

import (
	"fmt"
	"net"
	"suber/siface"
	"suber/utils"
)

type Server struct {
	//服务器名称
	Name string
	//服务器IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听端口
	Port int

	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务
	MsgHandler siface.IMsgHandler

	//该server的链接管理器
	ConnMgr siface.IConnManager
}

/*
func CallBackClient(conn *net.TCPConn,data []byte,cnt int)error{
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _,err := conn.Write(data[:cnt]);err != nil{
		fmt.Println("onn.Write(data[:cnt] err = )",err)
		return  errors.New("CallBackClient Error")
	}
	return nil
}
*/


//启动服务器
func (s *Server) Start(){
	fmt.Printf("[Suber] Server Name: %s, listenner at IP: %s,Port：%d\n",
		utils.GlobalConfigNow.Name,utils.GlobalConfigNow.Host,utils.GlobalConfigNow.TcpPort)
	fmt.Printf("[Suber] Version:%s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalConfigNow.Version,utils.GlobalConfigNow.MaxConn,utils.GlobalConfigNow.MaxPackageSize)


	go func() {
		//开启消息队列，及worker工作池
		s.MsgHandler.StartWorkerPool()
		// 1. 获取TCP ADD
		addr,err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil{
			fmt.Println("ResolveIPAddr err")
			return
		}
		//2. 监听服务器地址
		listener,err:=net.ListenTCP(s.IPVersion,addr)
		if err != nil{
			fmt.Println("ListenTCP(s.IPVersion,addr)")
			return
		}
		fmt.Println("start Suber server succ,",s.Name,"succ listenning")
		//3. 阻塞的等待客户端连接，处理客户端连接业务
		var cid uint32
		cid = 0
		for {
			conn,err:=listener.AcceptTCP()
			if err != nil{
				fmt.Println("listener.AcceptTCP() err=",err)
				continue
			}

			// 设置当前连接个数的判断，如果超过最大链接，则关闭此链接
			if s.ConnMgr.Len() >= utils.GlobalConfigNow.MaxConn{
				//TODO 给客户端相应一个超出最大连接的错误包
				fmt.Println("too many connetion....MaxConn = ",utils.GlobalConfigNow.MaxConn)
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn:=NewConnection(s,conn,cid,s.MsgHandler)
			cid++

			//启动当前链接的业务处理
			go dealConn.Start()
		}
	}()

}
//获取当前的server的链接管理器
func (s *Server)GetConnMgr() siface.IConnManager{
	return s.ConnMgr
}

//停止服务器
func (s *Server)Stop(){
	//TODO 将一些服务器资源、状态、或者开辟的链接信息进行停止或者回收
	fmt.Println("[server stop] suber server stop")
	s.ConnMgr.ClearConn()
}
//运行服务器
func (s *Server)Serve()  {
	//启动server服务器
	s.Start()
	//TODO 做一些启动之后的事情
	select {

	}
}
func (s *Server)AddRouter(msgID uint32,router siface.IRouter)  {
	s.MsgHandler.AddRouter(msgID,router)
	fmt.Println("add router success...")
}
/*
	初始化server模块方法
*/

func NewServer(name string) siface.IServer{
	s:=&Server{
		Name:utils.GlobalConfigNow.Name,
		IPVersion:"tcp4",
		IP:utils.GlobalConfigNow.Host,
		Port:utils.GlobalConfigNow.TcpPort,
		MsgHandler:NewMsghander(),
		ConnMgr:NewConnManager(),
	}
	return s
}