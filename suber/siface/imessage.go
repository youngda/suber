package siface

/*
	将请求的消息封装到一个Message中
*/
type IMessage interface{
	//获取消息ID
	GetMsgId() uint32
	//获取消息长度
	GetMsgLen() uint32
	//获取消息内容
	GetMsgData() []byte
	//设置消息ID
	SetMsgId(uint32)
	//设置消息内容
	SetMsgData([]byte)
	//设置消息长度
	SetMsgLen(uint32)
}