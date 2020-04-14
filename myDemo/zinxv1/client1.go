package main

import (
	"fmt"
	"time"
	"net"
	"zinx/znet"
	"io"
)

func main()  {
	fmt.Println("clint1 start")
	time.Sleep(1* time.Second)
	conn,err:=net.Dial("tcp","127.0.0.1:8091")
	if err!=nil{
		fmt.Println("client start err,exit!")
		return
	}
	for{
		//创建一个dp
		dp:=znet.NewDataPack()
		binaryMsg,err:=dp.Pack(znet.NewMsgPackage(1,[]byte("zinx1 v6 test send message")))
		if err!=nil{
			fmt.Println("client pack error",err)
			return
		}

		if _,err=conn.Write(binaryMsg);err!=nil{
			fmt.Println("client write error",err)
			return
		}

		//读取客户端的Msg Head 二进制8个字节
		binaryHaed:=make([]byte,dp.GetHeadLen())
		if _,err:=io.ReadFull(conn,binaryHaed);err!=nil{
			fmt.Println("client read msg head error",err)
			break
		}
		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msgHead,err:=dp.Unpack(binaryHaed)
		if err!=nil{
			fmt.Println("client unpack error",err)
			break
		}
		//根据dataLen 再次读取Data 放在msg.Data中
		var msgData []byte
		if msgHead.GetMsgLen()>0{
			msg:=msgHead.(*znet.Message)
			msgData=make([]byte,msgHead.GetMsgLen())
			if _,err:=io.ReadFull(conn,msgData);err!=nil{
				fmt.Println("read msg data error",err)
				break
			}
			//完整的一个消息已经读取完毕
			fmt.Println("-->Recv MsgID:",msg.Id,"msgLen:",msg.DataLen,"data:",string(msgData))
		}
		//cpu 阻塞
		time.Sleep(1* time.Second)
	}
}
