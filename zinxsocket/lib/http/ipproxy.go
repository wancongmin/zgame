package myhttp

import (
	"bangseller.com/lib/exception"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//正常返回，个数 = 给定个数 + 1
func GetProxy(num int) []string {
	return GetProxyByZima(num)
}

//芝麻代理
//http://www.zhimaruanjian.com/
//测试账号/密码 ranxl123 / ranxl123
const zimaGetIPUrl = "http://webapi.http.zhimacangku.com/getip?num=%d&type=1&pro=&city=0&yys=0&port=1&time=1&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions="

//提取返回数据格式
//IP:Port
//IP:Port
func GetProxyByZima(num int) []string {
	resp, err := http.Get(fmt.Sprintf(zimaGetIPUrl, num))
	exception.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	exception.CheckError(err)

	return strings.Split(string(body), "\r\n") //最后有一个空行，需要忽略
}
