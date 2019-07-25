package snet

import "suber/siface"

type Request struct {
	//已经和客户端建立好的链接
	conn siface.IConnection

	//得到客户端请求的数据
	msg siface.IMessage
}

//得到当前链接
func (r *Request)GetConnection() siface.IConnection{
	return r.conn
}
//得到请求的消息数据
func (r *Request)GetData() []byte{
	return r.msg.GetMsgData()
}

func (r *Request)GetDataID()uint32{
	return r.msg.GetMsgId()
}