package main

import (
	"fmt"
	"net"
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
		_,err:=conn.Write([]byte("helle Suberv0.2"))
		if err != nil{
			fmt.Println("conn.Write err = ",err)
			return
		}


		buf := make([]byte,512)
		cnt,err := conn.Read(buf)
		if err != nil{
			fmt.Println("conn.Read(buf) err = ",err)
			return
		}
		fmt.Printf("Server call back: %s , count: %d \n",buf[:cnt],cnt)

		time.Sleep(1*time.Second)
	}
}

/*
export GOROOT=/usr/local/go
export GOPATH=/mnt/g/go
export GOBIN=$GOPATH/bin
export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
*/