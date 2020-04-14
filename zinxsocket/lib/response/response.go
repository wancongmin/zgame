package resp

//API 接口统一返回值
type Response struct {
	Code    string      `json:"code"`    //成功返回信息，为空字符串，错误返回时，返回错误的代码，建议为有意义的字符串，如 LoginError,便于客户端多语言化
	Message string      `json:"message"` //提示信息，成功时
	Data    interface{} `json:"data"`    //成功时，返回结果，为JSON 字符串
}
