package mdb

import (
	"bangseller.com/lib/config"
	"bangseller.com/lib/exception"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

//部署方式：一个APP只对应一个数据库，一个卖家固定到一个或者多个APP，（但只能到一个数据库）
// 多个APP 可以到一个数据库

//数据库连接，直接通过全局变量访问
//改进路劲：继承 sqlx.DB , 然后继续用 Db 全局变量
var Db *sqlx.DB

const (
	database       = "Database"
	driverName     = "DriverName"
	dataSourceName = "DataSourceName"
	maxOpenConns   = "MaxOpenConns"
)

/**
初始打开数据库连接
{
	"Database":{"DriverName":"mysql","DataSourceName":"user:password@tcp(ip:port)/database","MaxOpenConns":10}
}
root:mysql@tcp(127.0.0.1:3306)/merp?charset=utf8&parseTime=true&loc=Local
关于连接串中的 parseTime=true, 这个会解析时间，包括时间和时区，为了方便，不要设置，所有时间全采用字符串的方式
*/
func InitDb() {
	dbm := config.GetMapConfig(database)
	var err error
	Db, err = sqlx.Open(dbm[driverName].(string), dbm[dataSourceName].(string))
	exception.CheckError(err)

	err = Db.Ping() //校验连接
	exception.CheckError(err)

	Db.SetMaxOpenConns(int(dbm[maxOpenConns].(float64))) // 最大打开连接数
	Db.SetMaxIdleConns(6)                                // 最大空闲连接数
	Db.SetConnMaxLifetime(60 * time.Minute)              //最大
	go ping()
	log.Println("数据库初始化成功")
}

//保证连接可用
func ping() {
	for {
		time.Sleep(5 * time.Minute)
		err := Db.Ping()
		if err != nil {
			fmt.Println(time.Now(), err)
		}
	}
}

//用于传递事务
type TxLink struct {
	Tx      *sqlx.Tx                   `json:"-"`
	MapStmt map[string]*sqlx.NamedStmt `json:"-"`
}

//释放资源
//在调用 tx.Tx = tx 后，调用 defer tx.Close()
// tx.Tx = tx
// defer tx.Close()
func (tx *TxLink) CloseTx() {
	//先 Close stmt,按入栈的反序
	if tx.MapStmt != nil {
		for _, stmt := range tx.MapStmt {
			if stmt != nil {
				stmt.Close()
			}
		}
	}
	//再 Rollback
	if tx.Tx != nil {
		tx.Tx.Rollback()
	}
}
