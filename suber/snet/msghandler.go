package snet

import (
	"fmt"
	"strconv"
	"suber/siface"
	"suber/utils"
)

/*
	消息处理模块的实现
*/
type MsgHander struct {
	Apis map[uint32]siface.IRouter
	//复制worker读取任务的消息队列
	TastQueue []chan siface.IRequest
	//业务工作池的工作数量
	WorkerPoolSize uint32

}

//初始化//创建MsgHander方法
func NewMsghander() *MsgHander {
	return &MsgHander{
		Apis:make(map[uint32] siface.IRouter),
		WorkerPoolSize:utils.GlobalConfigNow.WorkerPoolSize,//从全局配置中获取
		TastQueue:make([]chan siface.IRequest,utils.GlobalConfigNow.WorkerPoolSize),
	}
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

//启动一个Worker工作池
//开启工作池的动作只能启动一次工作池
func (mh *MsgHander)StartWorkerPool(){
	//根据workerPoolSize分别开启Worker,每个Worker用一个go承载
	for i:=0;i<int(mh.WorkerPoolSize) ;i++  {
		//一个worker被启动
		//当前的worker对应的channel消息队列 开辟空间|第0用第0个channel
		mh.TastQueue[i] =make(chan siface.IRequest,utils.GlobalConfigNow.MaxWorkerTaskLen)
		//启动当前的worker，阻塞等待消息从channel 中到来
		go mh.StartOneWorker(i,mh.TastQueue[i])

	}
}
//启动一个worker工作流程
func (mh *MsgHander)StartOneWorker(workID int,teskQueue chan siface.IRequest)  {
	fmt.Println("Worker ID = ",workID,"is started...")
	//不断的阻塞等待对应的消息队列的消息
	for{
		select {
			//如果有消息过来
			case request:=<-teskQueue:
				mh.DoMsgHandler(request)

		}
	}
}

//将消息交给TaskQueue，由Worker处理
func (mh *MsgHander)SendMsgToTaskQueue(request siface.IRequest)  {
	//将消息平均分配
	//根据客户端建立的ConnID来分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("add ConnID = ",request.GetConnection().GetConnID(),
		"request MsgID = ",request.GetDataID(),
		"to WorkerID = ",workerID,
		)
	// 将消息发送给对应的Worker的TaskQueue
	mh.TastQueue[workerID] <-request
}