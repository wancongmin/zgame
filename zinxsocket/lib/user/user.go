package user

import (
	"bangseller.com/lib/mdb"
	"bangseller.com/lib/exception"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type User struct {
	SellerId        int    `db:"seller_id" json:"seller_id"`
	UserId          int    `db:"user_id" json:"user_id"`
	UserCode        string `db:"user_code" json:"user_code"`
	UserName        string `db:"user_name" json:"user_name"`
	UserQq          string `db:"user_qq" json:"user_qq"`
	UserEmail       string `db:"user_email" json:"user_email"`
	UserPhone       string `db:"user_phone" json:"user_phone"`
	UserPwd         string `db:"user_pwd" json:"user_pwd"`
	UserStatus      int    `db:"user_status" json:"user_status"`
	DeptId          int    `db:"dept_id" json:"dept_id"`
	UserRoleId      int    `db:"user_role_id" json:"user_role_id"`
	CreateDate      string `db:"create_date" json:"create_date"`
	LastLoginDate   string `db:"last_login_date" json:"last_login_date"`
	LoginTimes      int    `db:"login_times" json:"login_times"`
	LoginErrorTimes int    `db:"login_error_times" json:"login_error_times"`
	Reffer          int    `db:"reffer" json:"reffer"`
}

//仅针对package可访问
var dc = map[string]string{
	mdb.CHILDRENALIA: "", //映射中一对多的表别名，多个逗号分开
	// sys_user
	"seller_id":         "u.seller_id",
	"user_id":           "u.user_id",
	"user_code":         "u.user_code",
	"user_name":         "u.user_name",
	"user_qq":           "u.user_qq",
	"user_email":        "u.user_email",
	"user_phone":        "u.user_phone",
	"user_pwd":          "u.user_pwd",
	"user_status":       "u.user_status",
	"dept_id":           "u.dept_id",
	"user_role_id":      "u.user_role_id",
	"create_date":       "u.create_date",
	"last_login_date":   "u.last_login_date",
	"login_times":       "u.login_times",
	"login_error_times": "u.login_error_times",
	"reffer":            "u.reffer",
}

//加密密码
//目前只是简单的MD5加密，下来升级加密算法，与user_id,create时间关联起来
func encryptPassword(user *User) string {
	h := md5.New()
	h.Write([]byte(user.UserPwd))
	return hex.EncodeToString(h.Sum(nil))
}

//验证登录密码
//暂时不处理identify，因go无session,所以 identify 将采用将 code 加密发到客户端，客户端传回进行校验
//校验成功返回 User 信息，否则返回 nil
func CheckLogin(useremail string, password string, identify string) User {
	user := GetUserByEmail(useremail)
	if len(user) == 0 {
		return User{}
	}
	fmt.Printf("%x", &user[0])
	return user[0]
}

//根据Email获取用户信息
func GetUserByEmail(useremail string) []User {
	user := []User{}
	err := mdb.Db.Select(&user, getUserByEmailSql, useremail)
	exception.CheckSqlError(err)
	return user
}
