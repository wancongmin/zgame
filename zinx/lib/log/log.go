package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// CreateDir 根据当前日期来创建文件夹
func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		err = os.MkdirAll(path, os.ModePerm) //0777也可以os.ModePerm
		if err != nil {                      //此处不能直接调用 exception.CheckError ,会造成循环引用
			log.Panic(err)
		}

		err = os.Chmod(path, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}
}

// 格式化文件路劲，路径格式为 logs/sellerid/date/path/filename
func FormatLogPath(sellerId int, path string, filename string) string {
	return fmt.Sprintf("logs/%d/%s/%s/%s", sellerId, time.Now().Format("20060102"), path, filename)
}

// 记录文件
func LogFile(filename string, data []byte) {
	index := strings.LastIndex(filename, "/")
	CreateDir(filename[:index]) //创建目录
	err := ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
}

//继承 Logger
type MyLogger struct {
	*log.Logger
	mu   sync.Mutex
	day  int
	file *os.File
}

//日志路径
var path string = "logs/"

// 重写Logger 基类，实现日志记录按小时
func (l *MyLogger) Output(calldepth int, s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	d := time.Now().Day()
	if l.day != d {
		l.day = d
		file, err := os.Create(path + time.Now().Format("20060102") + ".log")
		if err != nil {
			fmt.Println(err)
		} else {
			if l.file != nil { //新的文件创建成功，关闭前一文件
				l.file.Close()
				l.file = nil
			}
			l.file = file
			l.Logger.SetOutput(file)
		}
	}
	return l.Logger.Output(calldepth+1, s)
}

//为了使用原本的函数，变量也定义为 std
var std = NewLogger()

func NewLogger() MyLogger {
	pl := log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)
	return MyLogger{Logger: pl}
}

// 以下为从 go 源码中复制过来的，未作任何改变，仅仅是为了提供通过 package log 的访问
// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(2, s)
	panic(s)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, s string) error {
	return std.Output(calldepth+1, s) // +1 for this frame.
}
