package main

import (
	"fmt"
	"suber/siface"
	"suber/snet"
)
/*
	基于Suber 框架开发 服务器端应用程序
*/


//ping 测试 自定义路由

type PingRouter struct {
	snet.BaseRouter
}

//Test Prehandle
func (this *PingRouter)PreHandle(request siface.IRequest)  {
	fmt.Println("Call PreHandle... ")
	_,err:=request.GetConnection().GetTCPConnection().Write([]byte("Pre Ping\n"))
	if err != nil{
		fmt.Println("GetTCPConnection().Write Pre ping err = ",err)
	}
}
//Test handle
func (this *PingRouter)Handle(request siface.IRequest)  {
	fmt.Println("Call Handle... ")
	_,err:=request.GetConnection().GetTCPConnection().Write([]byte("Ping Ping\n"))
	if err != nil{
		fmt.Println("GetTCPConnection().Write Ping Ping err = ",err)
	}
}
//Test PostHandle
func (this *PingRouter)PostHandle(request siface.IRequest)  {
	fmt.Println("Call PostHandle... ")
	_,err:=request.GetConnection().GetTCPConnection().Write([]byte("Post Ping\n"))
	if err != nil{
		fmt.Println("GetTCPConnection().Write Post Ping err = ",err)
	}
}

func main()  {
	//创建一个server句柄
	s:= snet.NewServer("[suber v0.3]")

	//给当前suber框架添加一个自定义router
	s.AddRouter(&PingRouter{})
	//启动server
	s.Serve()
}