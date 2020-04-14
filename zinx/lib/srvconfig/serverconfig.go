package srvconfig

import (
	"bangseller.com/lib/mdb"
	"database/sql"
)

//服务端配置信息
//配置信息，将会同时保存到redis 和 MySQL中，获取的时候从Redis获取

type SysConfigure struct {
	SellerId    int            `db:"seller_id" json:"seller_id"`
	ConfigKey   string         `db:"config_key" json:"config_key"`
	ConfigValue sql.NullString `db:"config_value" json:"config_value"`
	LookType    sql.NullString `db:"look_type" json:"look_type"`
	DataType    sql.NullString `db:"data_type" json:"data_type"`
	Remark      sql.NullString `db:"remark" json:"remark"`
	InvalidDate sql.NullString `db:"invalid_date" json:"invalid_date"`
}

var dc = map[string]string{
	mdb.CHILDRENALIA: "", //映射中一对多的表别名，多个逗号分开
	// sys_configure
	"seller_id":    "cfg.seller_id",
	"config_key":   "cfg.config_key",
	"config_value": "cfg.config_value",
	"look_type":    "cfg.look_type",
	"data_type":    "cfg.data_type",
	"remark":       "cfg.remark",
	"invalid_date": "cfg.invalid_date",
}

func GetConfig(sellerId int, key string) string {
	v := ""
	mdb.Db.Get(&v, getConfigSql, sellerId, key)
	return v
}

func SetConfig() {

}
