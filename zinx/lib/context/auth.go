package context

import (
	"bangseller.com/lib/exception"
	msg "bangseller.com/lib/message"
	"bangseller.com/lib/redis"
	"bangseller.com/lib/user"
	"encoding/base64"
	"encoding/json"
	"time"
)

//校验权限接口
type IAuth interface {
	CheckAuth(c *Context) bool //使用中的校验
}

//鉴权实现
type Auth struct {
	SellerId     int    `json:"seller_id"`      //登录SellerId
	UserId       int    `json:"user_id"`        //登录用户ID,在UserId大于0的情况下，需要校验功能权限
	UserName     string `json:"user_name"`      //登录用户名
	UserEmail    string `json:"user_email"`     //登录用户邮件
	DeptId       int    `json:"dept_id"`        //登录用户部门
	UserRoleId   int    `json:"user_role_id"`   //登录用户角色
	UserRoleType string `json:"user_role_type"` //登录用户的基本角色
}

var (
	CookieTimeOut = 120 * time.Minute //cookie超时时间
)

const (
	sessionID     = "PHPSESSID" //在PHP存在的时候，就用PHP的SessionID
	sessionIDTime = "session-id-time"
	sessionToken  = "session-token"
	sellerToken   = "bangseller-x-token" //用于客户端卖家Token
	serverToken   = "bangseller-s-token" //服务调用Token，无需登录
)

//校验登录信息或者接口信息
//登录根据cookie来，在浏览器中，使用cookie是最方便的，不需要关心
//接口根据 Token 信息来
func CheckAuth(c *Context) bool {
	//优先校验Token
	if getAuthFromRedis(c) == true {
		return true
	}
	if c.sessionId == "" {
		return false
	}
	if c.isCookie == true {

	} else {
		return CheckToken(c)
	}
	return true
}

//首先从Redis中获取用户信息
func getAuthFromRedis(c *Context) bool {
	//exception.CheckError(err) 此处不能 CheckError ,在获取不到指定cookie的时候货抛出错误，只能判断 cookie 是否为 nil
	cookie, _ := c.Cookie(sessionID)
	if cookie != nil && cookie.Value != "" {
		c.sessionId = cookie.Value
		c.isCookie = true
	} else {
		c.isCookie = false
		c.sessionId = c.Header.Get(sellerToken)
		if c.sessionId == "" {
			return false
		} else {
			dest, _ := base64.StdEncoding.DecodeString(c.sessionId) //base64解密
			c.sessionId = string(dest)
		}
	}
	//从 Redis 获取数据
	if redis.GetStruct(c.sessionId, c.A) {
		if c.A.SellerId > 0 {
			//向后设置超时时间
			redis.Redis.Expire(c.sessionId, CookieTimeOut)
			return true
		} else {
			redis.Redis.Del(c.sessionId)
			return false
		}
	}
	return false
}

//将登录信息写入Redis
func SetAuthToRedis(c *Context) {
	redis.SetStruct(c.sessionId, c.A, CookieTimeOut)
}

//以登录方式鉴权，cookie
func CheckLogin(c *Context) bool {
	return true
}

//以Token方式鉴权
func CheckToken(c *Context) bool {
	c.A.SellerId = user.GetSellerId(c.sessionId)
	if c.A.SellerId == 0 {
		return false
	}
	SetAuthToRedis(c)
	return true
}

//暂时直接通过PHP传登录信息过来
func Login(c *Context) {
	c.sessionId = c.FormValue("sessionid")
	if c.sessionId == "" {
		exception.Throw(msg.AuthError)
		return
	}
	userinfo := c.FormValue("userinfo")
	err := json.Unmarshal([]byte(userinfo), c.A)
	exception.CheckError(err)

	if c.A.SellerId == 0 {
		exception.Throw(msg.AuthError)
		return
	}

	SetAuthToRedis(c)
	//cookie := &http.Cookie{Name: sessionID, Value: c.sessionId}
	//c.AddCookie(cookie)
}
