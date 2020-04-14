package channel

import (
	"bangseller.com/lib/config"
	"bangseller.com/lib/log"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type TaskParam struct {
	ClientId          int    `db:"client_id" json:"client_id"`                       //结点ID,用于标识这个是哪一个结点在处理
	PlanExecTime      string `db:"plan_exec_time" json:"plan_exec_time"`             //结点处理的时间,用于回滚标记
	PlanExecStatus    string `db:"plan_exec_status" json:"plan_exec_status"`         //计划执行状态
	PlanExecTimePoint string `db:"plan_exec_time_point" json:"plan_exec_time_point"` //计划执行时间点
}

const (
	TaskPending = "pending"
	TaskRunning = "running"
	TaskClose   = "close"
)

const (
	clientId    = "ClientId"
	workerCount = "WorkerCount"
)

//处理并发爬虫任务
//所有的任务都采用类似方式进行处理

//获取任务，返回需要执行任务的数组，必须是数组
type GetWorkFunc func(clientId int, workerCount int) interface{}

//处理任务
type DoWorkFunc func(workerId int, v interface{})

//服务进程ID
var ClientID int = config.GetIntConfig(clientId, 1)

//协程处理接口，需要实现这两个接口的
type ChannelInterface interface {
	GetWorkList(clientId int, workerCount int) interface{}
	DoWork(workerId int, v interface{})
}

//任务管理
type ChannelMgr struct {
	ch chan interface{} //Channel,任务队列

	Name            string      //协程名称, 用于从配置文件 config.json 中读取配置信息 "RT-Name":{"WorkerCount":50}
	WorkerCount     int         //开启任务数，同时取任务数
	GetWorkListFunc GetWorkFunc //获取任务的函数
	DoWorkFunc      DoWorkFunc  //执行任务的函数
	ChanInterface   *ChannelInterface
	errorCount      int //监控异常次数
}

//带错误处理，否则在循环中，出现异常，将会退出整个任务
func (c *ChannelMgr) getWorkWithHandleError() {
	defer HandleError()

	wl := c.GetWorkListFunc(ClientID, c.WorkerCount)
	v := reflect.ValueOf(wl) //将结果反射为数组，其他方式都不能做类型转换，只有等2.0的模板变量了
	l := v.Len()
	for i := 0; i < l; i++ {
		c.ch <- v.Index(i).Interface() //添加到任务队列
	}
	//取任务，等待1秒钟，解决时间+客户端id来控制并发冲突下取重复数据，这样也不会影响效率
	//或者采用唯一值方式，现在是
	time.Sleep(time.Second)
	runtime.Gosched() //释放资源
}

//获取需要执行的任务
func (c *ChannelMgr) getWork() {
	wgTask.Add(1)
	for {
		if isStop { //停止继续提交任务到队列，已经获取回来的，等待处理完成
			break
		}
		c.getWorkWithHandleError()
	}

	fmt.Printf("%s Routine End", c.Name)
	wgTask.Done()
}

//等待任务全部执行完成
func (c *ChannelMgr) waitWorkFinish() {
	for {
		if isStop && len(c.ch) == 0 {
			fmt.Println("关闭任务：", c.Name)
			close(c.ch) //关闭 channel
			break
		}
		time.Sleep(3 * time.Second)
	}
}

//开启任务
func (c *ChannelMgr) StartTask() {
	cf := config.GetMapConfig(c.Name)
	if cf == nil {
		panic(fmt.Sprintf("协程 %s 的配置信息不存在", c.Name))
	}
	c.WorkerCount = int(cf[workerCount].(float64))
	c.ch = make(chan interface{}, c.WorkerCount)
	go c.getWork()
	for i := 1; i <= c.WorkerCount; i++ {
		go c.doWork(i)
	}
	go c.waitWorkFinish()
}

// 带错误处理的执行任务
func (c *ChannelMgr) doWorkWihHandleError(workId int, t interface{}) {
	//调用任务执行函数
	c.DoWorkFunc(workId, t)
	runtime.Gosched()
}

// 执行任务
func (c *ChannelMgr) doWork(workId int) {
	wgWork.Add(1)

	for {
		//这种方式不行，因为所有的任务都被阻塞，根本不执行到这儿来，这样就释放不了
		//还是应该在主函数中控制
		//if isStop && len(c.ch) == 0 { //接收到停止信号后，所有任务执行完后退出
		//	fmt.Println("结束doWork", workId)
		//	break
		//}
		t, ok := <-c.ch  //存在任务，直接会返回任务；没有任务被阻塞（等待任务）
		if ok == false { //close ch 后，返回 false
			break
		}
		c.doWorkWihHandleError(workId, t)
	}
	fmt.Println("TaskWorker End:", workId)
	wgWork.Done()
}

//停止任务
var isStop = false

//GetWork 状态变量
var wgTask sync.WaitGroup

//DoWork 同步变量
var wgWork sync.WaitGroup

// SignalDeal 监控退出信号处理
//接收操作系统的信息
func MonitorSignal() {
	signalChan := make(chan os.Signal, 10)
	signal.Notify(signalChan,
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	s := <-signalChan
	isStop = true
	log.Println("接收到退出信号：", s)
}

//启动Rountine
//可以通过命令行参数启动任务
func StartMonitor() {
	MonitorSignal() //监控信号

	//等待提交任务结束
	wgTask.Wait()
	//等待所有任务执行结束
	wgWork.Wait()
}

//处理异常
func HandleError() {
	err := recover()
	if err != nil {
		log.Println(err)
	}
}
