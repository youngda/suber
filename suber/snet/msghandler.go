package snet

import (
	"fmt"
	"strconv"
	"suber/siface"
)

/*
	消息处理模块的实现
*/


//初始化//创建MsgHander方法
func NewMsghander() *MsgHander {
	return &MsgHander{
		Apis:make(map[uint32] siface.IRouter),
	}
}
type MsgHander struct {
	Apis map[uint32]siface.IRouter
}
//调度/执行对应的Router消息处理方法
func (mh *MsgHander) DoMsgHandler(request siface.IRequest){
	//从request 中找到msdID
	handler,ok:=mh.Apis[request.GetDataID()]
	if !ok{
		fmt.Println("api msgID = ",request.GetDataID(),"is not fount ,need register")
	}
	//根据MsgID 调度对应的Router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
//为消息添加具体的处理逻辑
func (mh *MsgHander) AddRouter(msgID uint32,router siface.IRouter){
	//判断当前的Msg绑定的API处理方法是否存在
	if _,ok:=mh.Apis[msgID];ok{
		//id已经注册
		panic("repeat api,msgID = "+ strconv.Itoa(int(msgID)))
	}
	//添加Msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID " , msgID," succ")
}