package main

import "suber/snet"
/*
	基于Suber 框架开发 服务器端应用程序
*/

func main()  {
	//创建一个server句柄
	s:= snet.NewServer("[suber v0.1]")
	//启动server
	s.Serve()
}