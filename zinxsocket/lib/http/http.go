package myhttp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"bangseller.com/lib/exception"
	"github.com/PuerkitoBio/goquery"
)

//请求参数结构
type RequestParam struct {
	Url         string
	Method      string // GET POST PATCH
	Body        io.Reader
	Req         *http.Request // Request 重新产生
	Client      *http.Client  // 多次访问，用同一Client, 每次Cookie都会自动传送 send 中处理
	ProxyUrl    string        // 记录方式，协议://ip:port
	Cookies     []*http.Cookie
	AgentMoblie int // 0 默认桌面; 1 Mobile
}

//电脑版Agent
var agents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3576.96 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3376.66 Safari/537.36",
}

//移动设备Agent
var agentsM = []string{
	"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Mobile Safari/537.36",
}

var agentLen int = len(agents)
var agentMLen int = len(agentsM)

//添加Request请求头信息
func SetHeader(req *http.Request, agentMobile int) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br") //加上此设置返回乱码，需要解压程序支持，所以去掉，后台不压缩
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if agentMobile == 1 {
		req.Header.Set("User-Agent", agentsM[rand.Intn(agentMLen)])
	} else {
		req.Header.Set("User-Agent", agents[rand.Intn(agentLen)])
	}
}

//一次性调用，多次调用或者带代理的请使用其他方法
func GetHtmlDoc(uri string) *goquery.Document {
	req, err := http.NewRequest("GET", uri, nil)
	exception.CheckError(err)

	SetHeader(req, 0)
	client := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	resp, err := client.Do(req)
	exception.CheckError(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	exception.CheckError(err)

	return doc
}

// 读取 Response 数据，如果压缩了，解压数据
func ReadData(response *http.Response) []byte {
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		exception.CheckError(err)
		data, err := ioutil.ReadAll(reader)
		exception.CheckError(err)
		return data
	default:
		data, err := ioutil.ReadAll(response.Body)
		exception.CheckError(err)
		return data
	}
}

/**
@param url string
*/
func GetByProxy(rq *RequestParam) *goquery.Document {
	req := NewRequest(rq)

	resp, err := rq.Client.Do(req)
	exception.CheckError(err)
	defer resp.Body.Close()

	data := ReadData(resp)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	exception.CheckError(err)

	return doc
}

// 设置代理服务器
// 1、代理模式
// 2、直接指定IP出口模式
func SetClientProxy(rq *RequestParam) {
	if strings.Index(rq.ProxyUrl, ":") >= 0 {
		//代理模式
		proxy, err := url.Parse(rq.ProxyUrl) // 直接 ip:port 格式是错误，需要 //ip:port (http)模式, socks5://ip:port
		exception.CheckError(err)
		rq.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
		return
	}
	rq.Client.Transport = &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			//本地地址  ipaddr是本地外网IP
			lAddr, err := net.ResolveTCPAddr(netw, rq.ProxyUrl+":0") // ":0" 是端口？
			if err != nil {
				return nil, err
			}
			//被请求的地址
			rAddr, err := net.ResolveTCPAddr(netw, addr)
			if err != nil {
				return nil, err
			}
			conn, err := net.DialTCP(netw, lAddr, rAddr)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(15 * time.Second))
			return conn, nil
		},
	}
}

// New 新的请求
func NewRequest(rq *RequestParam) *http.Request {
	req, err := http.NewRequest(rq.Method, rq.Url, rq.Body)
	exception.CheckError(err)

	if rq.Client != nil {
		req.Header = rq.Req.Header
		delete(req.Header, "Cookie") //删除
		rq.Req = req
		return req
	}

	rq.Req = req //保留下Request
	//初始化全新请求
	SetHeader(req, rq.AgentMoblie)
	jar, _ := cookiejar.New(nil)
	rq.Client = &http.Client{
		Jar:     jar, //设置这个后，在执行 client.Do 就会将返回的cookie保留下来
		Timeout: time.Duration(15 * time.Second),
	}
	if rq.Cookies != nil {
		//初始化Cookie，后续的Cookie直接就会返回到 Client 中
		rq.Client.Jar.SetCookies(req.URL, rq.Cookies)
	}

	if rq.ProxyUrl != "" {
		SetClientProxy(rq)
	}
	return req
}

//获取Cookie
func GetCookie(rq *RequestParam) *http.Response {
	req := NewRequest(rq)

	// Fetch Request
	resp, err := rq.Client.Do(req)
	exception.CheckError(err)
	return resp
}

func GetLocalIps() (ips []string) {
	addrs, err := net.InterfaceAddrs()
	exception.CheckError(err)

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
				fmt.Printf(`"%s",%s`, ipnet.IP.String(), "\n")
			}
		}
	}
	return ips
}
