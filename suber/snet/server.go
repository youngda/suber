package snet

import (
	"fmt"
	"net"
	"suber/siface"
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
	//当前的server添加一个router,server的注册的链接对应的处理业务
	Router siface.IRouter
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
	fmt.Printf("[Start] Server Listenner at IP: %s,Port %d, is starting\n",s.IP,s.Port)
	go func() {

		//1. 获取TCP ADD
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
			//将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn:=NewConnection(conn,cid,s.Router)
			cid++

			//启动当前链接的业务处理
			dealConn.Start()
		}
	}()

}
//停止服务器
func (s *Server)Stop(){
	//TODO 将一些服务器资源、状态、或者开辟的链接信息进行停止或者回收
}
//运行服务器
func (s *Server)Serve()  {
	//启动server服务器
	s.Start()
	//TODO 做一些启动之后的事情
	select {

	}
}
func (s *Server)AddRouter(router siface.IRouter)  {
	s.Router = router
	fmt.Println("add router success...")
}
/*
	初始化server模块方法
*/

func NewServer(name string) siface.IServer{
	s:=&Server{
		Name:name,
		IPVersion:"tcp4",
		IP:"0.0.0.0",
		Port:8999,
		Router:nil,
	}
	return s
}