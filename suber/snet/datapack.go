package snet

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"suber/siface"
	"suber/utils"
)

//封包拆包的具体模块
type DataPack struct {

}

//拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}
//获取包的长度的方法
func (d *DataPack) GetHeadLen()uint32{
	//DataLen Uint32(4字节) DataID unit32(4字节)
	return 8
}

//封包方法
/*
	dataLen|msgID|data
*/
func (d *DataPack) Pack(msg siface.IMessage) ([]byte,error){
	//创建一个存放bytes字节的缓冲
	dataBuff :=bytes.NewBuffer([]byte{})

	//将dataLen写到buff里
	if err:= binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgLen());err!=nil{
		return nil,err
	}
	//将MsgID写到buff里
	if err:= binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgId());err!=nil{
		return nil,err
	}
	//将data数据写到buff里
	if err:= binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgData());err!=nil{
		return nil,err
	}

	return dataBuff.Bytes(),nil
}
//拆包方法
/*
	将Head信息读出来，之后再根据Head信息里的data的长度，在进行一次读
*/
func (d *DataPack) UnPack(binData []byte)(siface.IMessage,error){
	//创建一个从输入二进制数据ioReader
	dataBuf:=bytes.NewReader(binData)

	//只解压head信息，得到MsgLen和MsgID
	msg := &Message{}
	//读MsgLen
	if err:= binary.Read(dataBuf,binary.LittleEndian,&msg.DataLen);err!=nil{
		return nil,err
	}
	//MsgID
	if err:= binary.Read(dataBuf,binary.LittleEndian,&msg.Id);err!=nil{
		return nil,err
	}
	//判断MsgLen是否已超过允许的最大包长度
	if utils.GlobalConfigNow.MaxPackageSize > 0 && msg.DataLen > utils.GlobalConfigNow.MaxPackageSize{
		return nil,errors.New("too large msg data recv!")
	}
	return msg,nil
}
