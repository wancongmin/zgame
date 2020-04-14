package mdb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bangseller.com/lib/exception"
	"bangseller.com/lib/fun"
	"github.com/jmoiron/sqlx"
)

/**
数据库相关的扩展操作，包括数据库连接、动态查询条件、分页等
*/

//存在一对多的多表联合查询列表原则:
//默认，按表头进行分页取数，需采用子查询，分组汇总
//如果存在表体的字段的查询条件或者排序，按明细方式进行查询及分页，不使用子查询，直接使用明细方式查询
//检查调用 CheckSqlChildren ,建议在字段对应map中通过 CHILDRENALIA 设置子表别名字符串，多个表以逗号分开

//sqlx 中,如果SQL语句中的字符串中含:，而非参数 ,需要使用 :: 进行转义
//建议返回数据，采用 Select ,通过结构返回，而不要使用 map 方式返回，根据SQL返回列表定义结构使用 devtool中的CreateSqlStruct

const (
	CHILDRENALIA = "_CHILDRENALIA"
	limit        = "limit"
	offset       = "offset"
	filter       = "F"
	pages        = "P"
)

//操作符号定义
var ops = map[string]string{
	"EQ":  "=",
	"GE":  ">=",
	"LE":  "<=",
	"GT":  ">",
	"LT":  "<",
	"LK":  "LIKE",
	"RK":  "LIKE",
	"AK":  "LIKE",
	"IN":  "IN",
	"NL":  "IS NULL",
	"NN":  "IS NOT NULL",
	"NEQ": "!=",
	"NE":  "!=",
}

//根据字段及值生成参数化查询条件
//返回生成的SQL语句，SQL中的参数通过 pm 字典返回，而不用再使用返回值返回了
func GetOpSQL(col string, v string, pm map[string]interface{}) string {
	query := ""
	op := ""

	p := strings.Replace(col, ".", "_", -1) //将小数点转为_作为参数的名称
	kArr := strings.Split(col, ".")
	l := len(kArr)
	if l <= 2 { //只有字段名字，默认操作符为 EQ(=)
		col = strings.Join(kArr, ".")
		op = "EQ"
	} else if l == 3 { //全部信息包括，其他的不做处理
		op = kArr[2]
		_, ok := ops[op]
		if !ok { //执行操作符不存在，忽略，不过此时可以回调函数进行处理
			return ""
		}
		col = kArr[0] + "." + kArr[1]
	}

	switch op {
	case "LK":
		if strings.Index(v, "%") < 0 {
			v = "%" + v + "%" //如果本身带有%号，不再添加%
		}
		break
	case "RK":
		if strings.Index(v, "%") < 0 {
			v = v + "%" //如果本身带有%号，不再添加%
		}
		break
	case "AK": //对于AK的语句，将值按空格拆分，进行全匹配
		vArr := strings.Split(v, " ")
		for i, v := range vArr {
			v = strings.TrimSpace(v)
			if v != "" {
				p1 := p + strconv.Itoa(i)
				query += " AND " + col + " LIKE :" + p1
				pm[p1] = "%" + v + "%"
			}
		}
		return query
	case "LE":
		//处理日期的LE，+1天，变成小于，主要是日期可能含有时间，未来可以将date 、datetime timestamp 区分开
		s := col[len(col)-5:]
		if s == "_date" || s == "_time" {
			t, err := time.Parse("2006-01-02", v)
			if err != nil { //日期转换错误，忽略
				return ""
			}
			v = t.Add(time.Hour * 24).Format("2006-01-02")
			op = "LT"
		}
		break
	case "IN": //IN中的值,不应该包含单引号，如果有，忽略，IN无法参数化
		v = strings.Replace(v, "'", "", -1)
		v = strings.Replace(v, ",", "','", -1)
		return " AND " + col + " IN ('" + v + "')"
	case "NL", "NN":
		return " AND " + col + " " + ops[op]
	}
	query = " AND " + col + " " + ops[op] + ":" + p
	pm[p] = v
	return query
}

//根据字段映射关系构造参数化查询条件
//这样前端就不需要看到表的别名，只需要 column[.op] 格式即可,省略 .op 代表 .EQ
//dc 为字段映射关系，映射为前端的值对应到具体数据库的表字段及别名，所以此时多个SQL语句的情况下，表别名保持一致
//dc = map[string]string{
//	"seller_id":"lmt.seller_id",
//	"listing_monitor_id":"lmt.listing_monitor_id",
//	"account_id":"lmt.account_id"
//	}
// 命名参数的值直接通过 m 返回
func GetWhereSQLByMap(m map[string]interface{}, dc map[string]string) string {
	if dc == nil { //dc为nil,按原本默认方式
		return GetWhereSQL(m)
	}

	fm := fun.GetSubMap(filter, m)
	query := bytes.Buffer{}

	for col, v := range fm {
		if v == "" {
			continue
		}
		//如果存在单引号，说明被人恶意纂改，忽略
		if strings.Index(col, "'") >= 0 {
			continue
		}

		//是否多字段,以逗号分开	//不再提供此种方式
		kArr := strings.Split(col, ".") //拆分字段和操作符
		if len(kArr) == 1 {
			kArr = append(kArr, "EQ") //默认等于操作符号
		}

		c, ok := dc[kArr[0]]
		if ok { //存在，用映射关系中的值，不存在，不创建SQL条件
			col = c + "." + kArr[1]
		} else {
			continue
		}

		query.WriteString(GetOpSQL(col, fmt.Sprint(v), m) + "\n")
	}

	return query.String()
}

/**
根据给定的map对象构造查询条件，对于动态条件(用户端输入的)，请一定调用本函数或者采用参数化，避免SQL注入
key 为 [alia.]column[.op]
生成命名参数
返回构造的SQL和对应参数的值map
给入Map 形如 {"F":{},"P":{},...},所以需要取出来
*/
func GetWhereSQL(m map[string]interface{}) string {
	fm := fun.GetSubMap(filter, m)
	query := bytes.Buffer{}

	for k, v := range fm {
		if v == "" {
			continue
		}
		//如果存在单引号，说明被人恶意纂改，忽略
		if strings.Index(k, "'") >= 0 {
			continue
		}
		//是否多字段,以逗号分开,不提供这种方式了，下来还是选择字段查询
		query.WriteString(GetOpSQL(k, fmt.Sprint(v), m) + "\n")
	}

	return query.String()
}

//分页结构定义
type Pager struct {
	Page      int    `json:"page"`       //指定页
	PageSize  int    `json:"page_size"`  //每页行数
	Offset    int    `json:"-"`          //Limit Offset
	PageCount int    `json:"page_count"` //页数
	Rows      int    `json:"rows"`       //总行数
	Col       string `json:"col"`        //排序字段
	Asc       string `json:"asc"`        //排序方式 ASC, DESC
}

//计算分页信息
func GetWherePageAndOrderBy(m map[string]interface{}, dc map[string]string) string {
	p := fun.GetSubMap(pages, m)

	pager := Pager{}
	fun.Map2Struct(p, &pager)

	pager.Offset = (pager.Page - 1) * pager.PageSize
	m[offset] = pager.Offset
	m[limit] = pager.PageSize

	//排序
	if v, ok := dc[pager.Col]; ok { //如果不存在字段
		pager.Col = v
	}
	orderby := pager.Col + " " + pager.Asc
	return orderby
}

//根据前端信息获取SQL语句需要的条件(where),分页(offset,limit),排序信息(orderby)
//m {"F":{},"P":{}}
//dc 字段映射关系，前端不显示表的别名
//返回值 where , orderby , 参数对应表
// 参数值直接通过 m 返回
func GetSqlInfo(m map[string]interface{}, dc map[string]string) (where string, orderby string) {
	where = GetWhereSQLByMap(m, dc)
	orderby = GetWherePageAndOrderBy(m, dc)
	return where, orderby
}

//根据map中的key值格式化SQL语句，替换形如{SQL}的模板部分，比如附加条件{where}、排序方式{orderby}等
//替换时区分大小写
func FormatSQLByMap(query string, m map[string]string) string {
	for k, v := range m {
		query = strings.Replace(query, "{"+k+"}", v, -1)
	}
	return query
}

//替换SQL语句中的{where}和{orderby}
//如果替换多的，请使用 FormatSQLByMap 函数，通过key值替换
func FormatSQL(query string, where string, orderby string) string {
	query = strings.Replace(query, "{where}", where, -1)
	if orderby != "" { //Order By 不为空时替换，为空时表示没有 {orderby}
		query = strings.Replace(query, "{orderby}", orderby, -1)
	}
	return query
}

//直接通过命名参数实现类似 sql.Select 功能
//arg 可以是struct, 也可以是 map[string]interface{}
func SelectNamed(dest interface{}, query string, arg interface{}) {
	stmt, err := Db.PrepareNamed(query)
	exception.CheckSqlError(err, 3)
	defer stmt.Close()

	err = stmt.Select(dest, arg)
	exception.CheckSqlError(err, 3)
}

//直接通过命名参数实现类似 sql.Get 功能
//arg 可以是struct, 也可以是 map[string]interface{}
func GetNamed(dest interface{}, query string, arg interface{}) {
	stmt, err := Db.PrepareNamed(query)
	exception.CheckSqlError(err, 3)
	defer stmt.Close()

	err = stmt.Get(dest, arg)
	exception.CheckSqlError(err, 3)
}

//检查一对多表是否存在表体的字段过滤和排序
//childrenTableAlias 子表别名,多个以逗号分开，不包含 .
func CheckSqlChildren(where string, orderby string, childrenTableAlias string) bool {
	for _, c := range strings.Split(childrenTableAlias, ",") {
		c = c + "."
		if strings.Contains(where, c) {
			return true
		}
		if strings.Contains(orderby, c) {
			return true
		}
	}
	return false
}

//Tx 默认没有通过命名参数直接调用 Select 方法
func SelectNamedTx(tx *sqlx.Tx, dest interface{}, query string, arg interface{}) {
	stmt, err := tx.PrepareNamed(query)
	exception.CheckSqlError(err, 3)
	defer stmt.Close()

	err = stmt.Select(dest, arg)
	exception.CheckSqlError(err, 3)
}

//Tx 默认没有通过命名参数直接调用 Get 方法
//get 在获取不到数据的时候也会返回异常，这种不好，所以少使用该方法
func GetNamedTx(tx *sqlx.Tx, dest interface{}, query string, arg interface{}) {
	stmt, err := tx.PrepareNamed(query)
	exception.CheckSqlError(err, 3)
	defer stmt.Close()

	err = stmt.Get(dest, arg)
	exception.CheckSqlError(err, 3)
}

// 替换参数打印出SQL语句
// 对于SQL中本身还有 :xx 的会有问题
func PrintSQL(query string, arg map[string]interface{}) {
	for k, v := range arg {
		sv := strings.Replace(fmt.Sprint(v), "'", "''", -1) //处理单引号
		query = strings.Replace(query, ":"+k, fmt.Sprintf("'%v'", sv), -1)
	}
	fmt.Println(query)
}
