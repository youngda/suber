package snet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只负责dataPack拆包封包单元测试
func TestDataPack_Pack(t *testing.T) {
	/*
		模拟服务器
	*/
	//创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("Server listen err:", err)
		return
	}
	//创建一个go 负责从客户端处理业务
	go func() {
		// 从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("listenner.Accept() err = ", err)
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				/*
					拆包的过程
				*/
				//定义一个拆包的对象
				dp := NewDataPack()
				for {
					//第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("io.ReadFull(conn,headData) err:", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("dp.UnPack(headData) err= ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//第二次从conn读，根据head中的MsgLen再读取data
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据dataLen的长度再次从io流中读取
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("io.ReadFull(conn,msg.Data) err:", err)
							return
						}

						//完整的消息已经读取完毕
						fmt.Printf("--->Recv MsgID:%d MsgLen:%d, MsgDada:%s\n", msg.Id, msg.DataLen, msg.Data)
					}

				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("conn,err := net.Dial")
		return
	}
	//创建一个封包对象dp
	dp := NewDataPack()

	//模拟粘包，两个包一起发

	//封装第一个包

	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'s', 'u', 'b', 'e', 'r'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("sendData,err := dp.Pack(msg1) err=", err)
		return
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 6,
		Data:    []byte{'h', 'i', 'g', 'i', 'r', 'l'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("sendData,err := dp.Pack(msg1) err=", err)
		return
	}

	//将两个包粘再一起
	sendData1  = append(sendData1,sendData2...)

	//一次性发送给服务端
	conn.Write(sendData1)


	//阻塞
	select {

	}
}
