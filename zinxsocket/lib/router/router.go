package router

import (
	"bangseller.com/lib/facebookgo/grace/gracehttp"
	"bangseller.com/lib/session"
	"fmt"
	"log"
	"net/http"
	"time"

	"bangseller.com/lib/config"
	"bangseller.com/lib/exception"
)

//处理函数类型
type HandleFunc func(c *session.Context)

//回调函数，处理权限验证等事情，直接抛异常退出
var callBackFunc HandleFunc

//路由表
//path : HandleFunc
//path 不需要后边的斜杠 / 如: /foo
var routerMap map[string]HandleFunc

//初始化Http服务，自动初始化数据库
//需要将配置文件设置好就行
//routers 路由表
//a 鉴权接口
func InitHttp(routers map[string]HandleFunc, callback HandleFunc) {
	defer func() { //异常捕获
		err := recover()
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}()

	routerMap = routers
	callBackFunc = callback

	//添加热更新发布，引用自Facebookgo,
	//github.com/facebookgo/grace/gracedemo
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainRouter)
	err := gracehttp.Serve(&http.Server{
		Addr:    ":" + config.GetConfig("port", "80"),
		Handler: mux,
	})
	if err != nil {
		fmt.Println(err)
	}
}

//主路由函数
//1、处理路由的分配
//2、登录验证：1)密码验证；2)Token 验证；3)Session验证
func mainRouter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		//浏览器因素，访问时都会来访问一次该链接,取消
		//在发布后，可以在 nginx中配置，不访问进来
		return
	}
	t := time.Now()
	//统一异常处理
	defer exception.HandleError(w)

	//登录校验，权限校验
	//https://github.com/dgrijalva/jwt-go 可以采用 JWT 方式进行
	//Session cookie ID
	c := &session.Context{
		Request: r,
		W:       w,
		A:       &session.Auth{},
	} //继承赋值方式为用父 struct 的名字

	c.AccessCountrolAllow() //统一调用

	path := r.URL.Path

	handleFunc, ok := routerMap[path]
	if !ok {
		//路由不存在，404错误
		c.Fail("404")
		return
	}

	//调用函数，如果是 form-multidata,自己调用解析函数	r.ParseMultipartForm()和r.MultipartForm获取值
	r.ParseForm() //解析参数

	//回调，用于校验用户身份信息，权限统一控制等
	if callBackFunc != nil {
		callBackFunc(c) //抛异常结束
	}
	handleFunc(c)
	fmt.Printf("%s\t%s\t%v\n", time.Now().Format("2006-01-02 15:04:05"), c.URL, time.Now().Sub(t))
	//	记录时间日志，记录处理时间长的功能，便于后续优化
}

// TestRouter 路由测试
func TestRouter(c *session.Context) {
	c.Success("路由加载成功:" + c.URL.Path)
}
