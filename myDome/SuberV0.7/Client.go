package main

import (
	"fmt"
	"io"
	"net"
	"suber/snet"
	"time"
)

/*
	模拟客户端
*/
func main()  {
	fmt.Println("Client start...")
	time.Sleep(1*time.Second)
	//1 直接链接远程服务器，得到一个conn链接

	conn,err:= net.Dial("tcp","127.0.0.1:7778")
	if err != nil{
		fmt.Println("net.Dial err = ",err)
		return
	}
	//链接调用Write 写数据
	for{
		//发送封包的Msg格式消息
		dp:=snet.NewDataPack()

		binMsg,err :=dp.Pack(snet.NewMsgPackage(0,[]byte("suber0.7")))
		if err != nil{
			fmt.Println("dp.Pack(snet.NewMsgPackage err = ",err)
			return
		}
		if _,err := conn.Write(binMsg);err != nil{
			fmt.Println("conn.Write(binMsg) err = ",err)
			return
		}

		/*
			服务器回复消息
		*/
		//先读取流中的head部分，得到ID和dataLen
		binHead := make([]byte,dp.GetHeadLen())
		if _,err:= io.ReadFull(conn,binHead);err!= nil{
			fmt.Println("io.ReadFull(conn,binHead) err = ",err)
			break
		}
		//将二进制的head拆包到Msg结构体中
		msgHead,err:= dp.UnPack(binHead)
		if err!= nil{
			fmt.Println("msgHead,err:= dp.UnPack(binHead) err = ",err)
			break
		}
		if msgHead.GetMsgLen()>0 {
			msg:=msgHead.(*snet.Message)
			msg.Data = make([]byte,msg.GetMsgLen())

			if _,err:=io.ReadFull(conn,msg.Data);err!= nil{
				fmt.Println("io.ReadFull(conn,msg.Data) err = ",err)
				return
			}
			fmt.Println("-->Recv Server Msg: ID = ",msg.Id,"Len = ",msg.DataLen,"data = ",string(msg.Data))
		}
		//再根据DataLen进行第二次读取
		//cpu阻塞
		time.Sleep(1*time.Second)
	}
}

/*
export GOROOT=/usr/local/go
export GOPATH=/mnt/g/go
export GOBIN=$GOPATH/bin
export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
*/