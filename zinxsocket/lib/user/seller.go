package user

import (
	"bangseller.com/lib/mdb"
	"bangseller.com/lib/exception"
)

//根据Token信息获取SellerId
func GetSellerId(token string) int {
	sellerId := 0
	err := mdb.Db.Get(&sellerId, getSellerIdByTokenSql, token)
	exception.CheckSqlError(err)
	if sellerId == 0 {
		exception.Throw("TokenError")
	}
	return sellerId
}
