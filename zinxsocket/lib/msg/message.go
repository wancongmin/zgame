package msg

//Message 常量定义，便于统一进行提示信息的维护
//在每个实现的地方，直接调用即可
//定义中按 package 放在一块
var M = map[string]string{
	// 404
	"404":     "您访问的页面不存在",
	Exception: "Sorry,系统执行异常,请重试，如果问题继续存在，请联系BangSeller客服解决",

	// Auth
	AuthError: "身份校验错误",

	// user/seller
	TokenError: "Token信息错误",

	StockTestLimit: "超过2个测试的上限",

	OutOfBalance:          "您的资金余额不足以支付本次消费金额！请充值后继续！",
	PageViewsError:        "Page Views 数据格式存在错误，请联系BangSeller客服",
	AccountEmailAuthError: "店铺邮件发送未授权",
	AccountApiAuthError:   "店铺未授权",
}

//提示字符串常量定义
const StockTestLimit = "StockTestLimit"
const OutOfBalance = "OutOfBalance"
const TokenError = "TokenError"
const AuthError = "AuthError"
const PageViewsError = "PageViewsError"
const Exception = "Exception"
const AccountEmailAuthError = "AccountEmailAuthError"
const AccountApiAuthError = "AccountApiAuthError"
