package main
import (
	"zinx/znet"
	"fmt"
	"zinx/ziface"
	"zinx/lib/config"
)
//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest){
	fmt.Println("Call PingRouter Handle..")
	//先读取客户端数据再回写
	fmt.Println("recv from client msgID=",request.GetMsgId(),",data=",string(request.GetData()))
	err:=request.GetConnection().SendMsg(200,[]byte("ping...ping..."))
	if err!=nil{
		fmt.Println("Handle err",err)
	}
}

type HolleRouter struct {
	znet.BaseRouter
}

func (this *HolleRouter) Handle(request ziface.IRequest){
	fmt.Println("Call HolleRouter Handle..")
	//先读取客户端数据再回写
	fmt.Println("recv from client msgID=",request.GetMsgId(),",data=",string(request.GetData()))
	err:=request.GetConnection().SendMsg(201,[]byte("hello...hello..."))
	if err!=nil{
		fmt.Println("Handle err",err)
	}
}
//创建链接之后执行的钩子函数
func DoConnectionBegin(conn ziface.Iconnection)  {
	fmt.Println("====>DoConnection is Call")
	if err:=conn.SendMsg(202,[]byte("DoConnection Beagin"));err!=nil{
		fmt.Println(err)
	}
	//链接之前设置一些属性
	fmt.Println("Set Property....")
	conn.SetProperty("name","wancongmin")
	conn.SetProperty("home","wuhan")
}

//链接断开执行的钩子函数
func DoConnectionLost(conn ziface.Iconnection)  {
	fmt.Println("====>DoConnectionLost is Call")
	fmt.Println("====>conn ID =",conn.GetConnID())
	//获取链接属性
	if val,err:=conn.GetProperty("name");err==nil{
		fmt.Println("name",val)
	}
	if val,err:=conn.GetProperty("home");err==nil{
		fmt.Println("home",val)
	}
}


func main(){
	//创建server句柄，使用zinx的api
	config.InitConfig()                   //初始化本地配置文件 config.json
	var Platform string = config.GetConfig("Platform","Linux/centos7")
	fmt.Println(Platform)
	s:=znet.NewServer("zinx06")
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HolleRouter{})

	//注册连接的Hook钩子函数
	s.SetConnStart(DoConnectionBegin)
	s.SetConnStop(DoConnectionLost)
	//启动Server
	s.Server()
}
