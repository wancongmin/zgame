package exception

import (
	"bangseller.com/lib/log"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"bangseller.com/lib/msg"
	resp "bangseller.com/lib/response"
)

// 显示异常
//以JSON方式返回结果
//在Controller层的http请求的入口程序的第一句话，通过 defer HandleError(w) 调用
func HandleError(w http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}

	var res resp.Response
	switch err.(type) {
	case resp.Response: //正常业务错误提示，直接返回信息
		res = err.(resp.Response)
	default:
		res = resp.Response{Code: msg.Exception, Message: fmt.Sprint(err)}
		//记录不是通过 Check...Error抛出的错误信息，层次 5 是比较合适的，基本能够定位到引起错误的地方
		log.Output(5, fmt.Sprint(err))
		log.Println(strings.Split(string(debug.Stack()), "main.go")[0])
	}
	data, err := json.Marshal(&res)
	if err != nil {
		fmt.Fprint(w, "JSON异常:", res.Data, err)
	} else {
		fmt.Fprintln(w, string(data))
	}
}

//检查异常，对于异常直接抛出，结束程序的运行
//如果是错误，需要通过正常途径返回
//只要函数有返回error的，不能不做错误的处理，否则，将会出现严重问题
func CheckError(err error) {
	if err != nil {
		//后台记录日志
		log.Output(2, fmt.Sprintln(err)) //此处直接调用 log.Println 打印出来的文件及行号为本文件的
		//抛出错误到调用程序
		res := resp.Response{Code: msg.Exception, Message: fmt.Sprintln(err)}
		panic(res)
	}
}

//在调用 db 的 Get 方法时，没有返回信息，会抛出错误，这种在程序处理中本来就是检查返回值，所以忽略错误的抛出
const getNoRowsError = "sql: no rows in result set"

//专门用于检查SQL错误
func CheckSqlError(err error, calldepth ...int) {
	if err != nil && err != sql.ErrNoRows { //Get 无结果返回
		if len(calldepth) == 0 {
			log.Output(2, fmt.Sprintln(err))
		} else {
			log.Output(calldepth[0], fmt.Sprintln(err))
		}
		res := resp.Response{Code: msg.Exception, Message: fmt.Sprintln(err)}
		panic(res)
	}
}

//检查在事务中的错误返回，如果存在错误，回滚事务
// 不再使用这种方式来 Rollback 事务，如果出现任务异常，将会出现无法 rollback
// 采用 开始事务后，直接 defer tx.Rollback()  这样不管如何都会调用 rollback, Commit 后调用也没关系
//func CheckTxError(err error, tx *sqlx.Tx, calldepth ...int) {
//	if err != nil && err != sql.ErrNoRows {
//		tx.Rollback()
//		if len(calldepth) == 0 {
//			log.Output(2, fmt.Sprintln(err))
//		} else {
//			log.Output(calldepth[0], fmt.Sprintln(err))
//		}
//		res := resp.Response{Code: msg.Exception, Message: fmt.Sprintln(err)}
//		panic(res)
//	}
//}

//通过panic抛出异常，只要是程序中遇到错误需要结束整个运行过程的，就调用此方法抛出异常，最后统一由 HandleError 处理
//msgkey 定义在 message 中
func Throw(msgkey string) {
	res := resp.Response{Code: msgkey, Message: msg.M[msgkey]}
	log.Output(2, fmt.Sprintln(res.Message))
	panic(res)
}

//记录错误日志
func LogError() {
	if err := recover(); err != nil {
		log.Output(2, fmt.Sprintln(err))
	}
}
