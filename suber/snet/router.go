package snet

import "suber/siface"

//实现router时，先嵌入BaseRouter基类，然后根据需要对这个基类方法进行重写
type BaseRouter struct {

}

//在处理conn业务之前的钩子方法
func(b *BaseRouter) PreHandle(request siface.IRequest){

}
//在处理conn业务之中的主方法
func(b *BaseRouter) Handle(request siface.IRequest){

}
//在处理conn业务之前的钩子方法
func(b *BaseRouter) PostHandle(request siface.IRequest){

}