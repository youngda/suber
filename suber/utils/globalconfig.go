package utils

import (
	"encoding/json"
	"io/ioutil"
	"suber/siface"
)

/*
	存储suber 框架的全局参数
	一些参数通过suber.json有用户进行配置
*/

type GloblaConfig struct {
	/*
		Server
	*/
	TcpServer siface.IServer //当前suber 全局的Server对象
	Host string				 //服务器监听的IP
	TcpPort int              //服务器监听的端口
	Name string              //服务器监听的名称

	/*
		Suber
	*/
	Version string//Suber版本号
	MaxConn int//允许的最大链接数
	MaxPackageSize uint32 //数据包的最大值
}

/*
	定义一个全局的对外GlobalConfig
*/
var GlobalConfigNow *GloblaConfig


/*
	从suber.json加载用户自定义设置参数
*/
func (g *GloblaConfig)Relaod()  {
	data,err:=ioutil.ReadFile("conf/suber.json")
	if err != nil{
		panic(err)
	}
	err = json.Unmarshal(data,&GlobalConfigNow)
	if err != nil{
		panic(err)
	}
}
/*
	提供一个init方法，初始化当前的全局对象
*/
func init()  {
	//如果配置文件没有加载，默认的值
	GlobalConfigNow = &GloblaConfig{
		Name:"SuberAPP",
		Version:"v0.4",
		TcpPort:8999,
		Host:"0.0.0.0",
		MaxConn:1000,
		MaxPackageSize:4096,
	}

	//尝试从配置文件中加载
	GlobalConfigNow.Relaod()
}