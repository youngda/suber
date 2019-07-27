package snet

import (
	"errors"
	"fmt"
	"suber/siface"
	"sync"
)

/*
	链接管理模块
*/
type ConnManager struct {
	//管理链接集合
	connections map[uint32] siface.IConnection
	//保护链接的读写锁
	connLock sync.RWMutex
}

//创建当前链接的方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections:make(map[uint32] siface.IConnection),
	}
}

//添加链接
func (cM *ConnManager) Add(conn siface.IConnection){
	//保护共享资源map ,加写锁
	cM.connLock.Lock()
	defer cM.connLock.Unlock()

	//将conn加入到ConnManager中
	cM.connections[conn.GetConnID()] = conn
	fmt.Println("conn id = ",conn.GetConnID(),"connection add to Connmanager successfully...conn num = ", len(cM.connections))

}
//删除链接
func (cM *ConnManager) Remove(conn siface.IConnection){
	//保护共享资源map ,加写锁
	cM.connLock.Lock()
	defer cM.connLock.Unlock()
	//删除链接器conn
	delete(cM.connections,conn.GetConnID())
	fmt.Println("conn id = ",conn.GetConnID(),"connection remove to Connmanager successfully...conn num = ", cM.Len())
}

//得到当前连接总数
func (cM *ConnManager) Len()int{
	return len(cM.connections)
}


//得到链接
func (cM *ConnManager) Get(connID uint32)(siface.IConnection,error){
	//保护共享资源map ,加读锁
	cM.connLock.RLock()
	defer cM.connLock.RUnlock()

	if conn,ok:= cM.connections[connID];ok{
		//找到
		return conn,nil
	}else {
		return nil,errors.New("connection not found")
	}
}


//清除并终止所有链接
func (cM *ConnManager) ClearConn(){
	//保护共享资源map ,加写锁
	cM.connLock.Lock()
	defer cM.connLock.Unlock()

	//删除conn,并停止conn的工作
	for conID,conn := range cM.connections{
		//停止
		conn.Stop()
		//删除
		delete(cM.connections,conID)
	}
	fmt.Println("clear all connections successful ! conn num = ",cM.Len())
}