package hothttp

// /**************************************************************************************
// Code Description    : 日志
// Code Vesion         :
// 					|------------------------------------------------------------|
// 						  Version    					Editor            Time
// 							1.0        					yuansudong        2016.4.12
// 					|------------------------------------------------------------|
// Version Description	:
//                     |------------------------------------------------------------|
// 						  Version
// 							1.0
// 								 ....
// 					|------------------------------------------------------------|
// ***************************************************************************************/
// package hothttp

// import (
// 	"fmt"
// 	"os"
// 	"runtime"
// 	"sync"

// 	"bangseller.com/lib/base/util"
// )

// var (
// 	dbg string
// 	inf string
// 	err string
// )

// const (
// 	// debugLevel Debug级别
// 	debugLevel = 0
// 	// infoLevel Info 级别
// 	infoLevel = 1
// 	// errLevel Error 级别
// 	errLevel = 2
// )

// // 文件句柄
// var fileHandle = os.Stdout

// func init() {
// 	inf = "[ INFO ]"
// 	dbg = "[ DEBUG ]"
// 	err = "[ ERROR ]"
// }

// // log
// type log struct {
// 	level      int
// 	fileHandle *os.File
// 	mutex      sync.Mutex
// }

// // newLog 用于
// func newLog() *log {
// 	pInst := new(log)
// 	pInst.level = debugLevel
// 	pInst.fileHandle = os.Stdout
// 	return pInst
// }

// func (l *log) write() {

// }

// // debug 用于打印调试信息
// f unc (l *log) debug(stack, format string, args ...interface{}) {
// 	if debugLevel >= l.level {
// 		args1 := []interface{}{
// 			util.GetDate(),
// 			"   ",
// 			dbg,
// 			"   ",
// 			stack,
// 			"   ",
// 			fmt.Sprintf(format, args...),
// 		}
// 		fmt.Fprintln(fileHandle, args1...)
// 	}
// }

// // Info 用于打印调试信息
// func (l *log) info(stack, format string, args ...interface{}) {
// 	if LogInfoLevel >= l.level {
// 		args1 := []interface{}{
// 			util.GetDate(),
// 			"   ",
// 			inf,
// 			"   ",
// 			stack,
// 			"   ",
// 			fmt.Sprintf(format, args...),
// 		}
// 		fmt.Println(
// 			args1...,
// 		)
// 	}
// }

// // Info 用于打印调试信息
// func (l *log) error(stack, format string, args ...interface{}) {
// 	if LogErrLevel >= l.level {
// 		args1 := []interface{}{
// 			util.GetDate(),
// 			"   ",
// 			err,
// 			"   ",
// 			stack,
// 			"   ",
// 			fmt.Sprintf(format, args...),
// 		}
// 		fmt.Println(
// 			args1...,
// 		)
// 	}
// }

// // Debug 用于打印Debug类型的消息
// func Debug(format string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	logHandle.debug(fmt.Sprintf("%s:%d", file, line), format, args...)
// }

// // Info 用于打印Info级别的消息
// func Info(format string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	logHandle.info(fmt.Sprintf("%s:%d", file, line), format, args...)
// }

// // Error 用于打印错误级别的消息
// func Error(format string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	logHandle.error(fmt.Sprintf("%s:%d", file, line), format, args...)
// }
