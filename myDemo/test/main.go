package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)
type IM interface {
	//获取消息的Id
	SetId()
	GetId()
}

type Ms struct {
	Id  uint32   //消息的id
}

//获取消息的Id
func (m *Ms) SetId(id uint32) {
	m.Id=id
}

func (m *Ms) GetId()uint32 {
	return m.Id
}

func test() ziface.IMessage {
	msg:=&znet.Message{}
	fmt.Println(msg.GetMsgLen())
	msg.Id=15
	return msg
}


func main() {

	m:=test()
	fmt.Printf("%+v",m)
	fmt.Println(m.Id)
}