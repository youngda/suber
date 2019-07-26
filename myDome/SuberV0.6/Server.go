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


//Test handle
func (this *PingRouter)Handle(request siface.IRequest)  {
	fmt.Println("Call Handle... ")
	//先读取客户端的数据，再回写ping
	fmt.Println("recv from client MsgId = ",request.GetDataID(),"data = ",string(request.GetData()))
	if err:=request.GetConnection().SendMsg(1,[]byte("ping...ping"));err!=nil{
		fmt.Println("request.GetConnection().SendMsg err",err)
	}
}


//HelloHandler 测试 自定义路由
type HelloHandler struct {
	snet.BaseRouter
}
//Test handle
func (this *HelloHandler)Handle(request siface.IRequest)  {
	fmt.Println("Call HelloHandle... ")
	//先读取客户端的数据，再回写ping
	fmt.Println("recv from client MsgId = ",request.GetDataID(),"data = ",string(request.GetData()))
	if err:=request.GetConnection().SendMsg(201,[]byte("welcome suber 0.6"));err!=nil{
		fmt.Println("request.GetConnection().SendMsg err",err)
	}
}


func main()  {
	//创建一个server句柄
	s:= snet.NewServer("[suber v0.6]")

	//给当前suber框架添加自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloHandler{})
	//启动server
	s.Serve()
}