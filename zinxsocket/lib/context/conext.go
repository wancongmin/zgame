package context

import (
	tc "bangseller.com/columns"
	"bangseller.com/lib/exception"
	"bangseller.com/lib/message"
	"bangseller.com/lib/response"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//将参数组合，适于未来的扩展
type Context struct {
	*http.Request //继承
	W             http.ResponseWriter

	A *Auth //鉴权信息，在实现 Auth的Check接口中处理

	sessionId string //记录用户唯一的SessionID凭证，包括token 或者 cookieID
	isCookie  bool   //是否来自于Cookie

	M interface{} //用于通过Context传递信息
}

//允许跨域访问
func (c *Context) AccessCountrolAllow() {
	w := c.W
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "bangseller-x-token,content-Type")
	w.Header().Set("content-type", "application/json;charset=utf-8") //返回数据格式是json
}

//提供直接通过 c 访问返回值
//不再使用 sql.null... 类型，通过指针解决更方便
func (c *Context) Success(data interface{}) {
	res := resp.Response{Code: "0", Message: "success", Data: data}
	c.AccessCountrolAllow() //允许跨域访问，解决Ajax调用返回到 success 回调中
	d, _ := json.Marshal(res)
	fmt.Fprint(c.W, string(d))
}

//访问错误返回，其后无执行代码的时候调用该方法
//在controller中遇到问题，直接调用该函数后return 结束程序，不要抛出异常
//msgkey 请在 message 进行维护
func (c *Context) Fail(msgkey string) {
	res := &resp.Response{Code: msgkey, Message: msg.M[msgkey]}
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(c.W, "JSON异常:", res.Data, err)
	} else {
		c.AccessCountrolAllow()
		fmt.Fprint(c.W, string(data))
	}
}

//获取int 参数
func (c *Context) FormIntValue(key string) int {
	s := strings.TrimSpace(c.FormValue(key))
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	exception.CheckError(err)
	return i
}

//获取Body 中通过 raw 方式提交的数据 包括application/json，text, Text/Plan 等等
func (c *Context) GetRawData() []byte {
	data, err := ioutil.ReadAll(c.Body)
	exception.CheckError(err)

	return data
}

//获取 application/json 数据
func (c *Context) GetMap() map[string]interface{} {
	m := make(map[string]interface{})
	err := json.NewDecoder(c.Body).Decode(&m)
	exception.CheckError(err)
	return m
}

//获取 application/json 数据
//直接设置SellerID 和 UserID
func (c *Context) GetMapWithSeller() map[string]interface{} {
	m := make(map[string]interface{})
	err := json.NewDecoder(c.Body).Decode(&m)
	exception.CheckError(err)

	m[tc.SellerId] = c.A.SellerId
	m[tc.UserId] = c.A.UserId
	return m
}

//获取 application/json 数据,结构体
func (c *Context) GetJsonStruct(s interface{}) {
	err := json.NewDecoder(c.Body).Decode(s) //也可以使用这种模式获取JSON结构
	exception.CheckError(err)

	//data, err := ioutil.ReadAll(c.Body) //被读一次后，再次读数据就读不到
	//exception.CheckError(err)
	//
	//err = json.Unmarshal(data, s)
	//exception.CheckError(err)
}
