package fun

import (
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bangseller.com/lib/exception"
)

/**
补充常用函数，golang自身未提供的函数
*/

const (
	Ymd    = "2006-01-02"
	YmdHis = "2006-01-02 15:04:05"
)

/**
round
*/
func Round(f float64, dec int) float64 {
	d10 := float64(math.Pow10(dec))
	return math.Round(f*d10) / d10
}

//日期转换，封装异常处理
func ParseDate(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	exception.CheckError(err)
	return t
}

/**
基于UTF8的字符串截取
*/
func SubString(s string, start int, length int) string {
	rs := []rune(s)
	return string(rs[start:(start + length)])
}

//String to Time
func StringtoDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	exception.CheckError(err)
	return t
}

//默认时区,设置为数字
var DefaultZone, _ = time.LoadLocation("Asia/Chongqing") //3600秒
const Minute2Second = 60                                 //每小时秒数
const Second2Ns = 1000000000                             //秒转化为Ns
func Now() time.Time {
	//使用这种方式是因为另外一种需要库的支持
	// 在windows平台上默认没有，如果发布到客户端，就会有问题
	return time.Now().In(DefaultZone)
}

//国家对应时区,支持夏令时
//参考 http://unicode.org/cldr/data/common/supplemental/windowsZones.xml
// 在go的源码中 go/time/zoneinfo_abbrs_windows.go 中有定义
var CountryTimeZone = map[string]string{
	"CN": "Asia/Chongqing",      //中国
	"JP": "Asia/Tokyo",          //JST,JST
	"US": "America/Los_Angeles", //PST,PDT
	"CA": "America/Los_Angeles", //PST,PDT
	"UK": "Europe/London",       //GMT,BST
	"FR": "Europe/Berlin",       //CET, CEST
	"DE": "Europe/Berlin",       //CET, CEST
	"IT": "Europe/Berlin",       //CET, CEST
	"ES": "Europe/Berlin",       //CET, CEST
}

func GetCountryZone(country string) *time.Location {
	zone, err := time.LoadLocation(CountryTimeZone[country])
	exception.CheckError(err)
	return zone
}

//获取当前时间的下一时间点，转为中国时间
//times 为时间 Hour:Minute 字符串，多个逗号分开, Hour 和 Minute 都为 2位模式 如: 08:30,18:05
func GetNextTime(country string, times string) time.Time {
	cz := GetCountryZone(country)
	t := Now().In(cz)

	hm := t.Format("15:04") //时分
	hs := strings.Split(times, ",")
	for _, v := range hs {
		if hm < v {
			dt := t.Format(Ymd) + " " + v + ":00"
			t, _ = time.ParseInLocation(YmdHis, dt, cz)
			return t.In(DefaultZone)
		}
	}
	//当前时间大于给定节点,下一天的第一个节点的时间
	dt := t.Add(24*time.Hour).Format(Ymd) + " " + hs[0] + ":00"
	t, _ = time.ParseInLocation(YmdHis, dt, cz)
	return t.In(DefaultZone)
}

//在给定时间的基础上加上时分秒
//该方法在处理时间的时候不改变原有时间的时区
func AddTime(t time.Time, hour int, minute int, second int) time.Time {
	if hour == 0 && minute == 0 && second == 0 {
		return t
	}
	var ns int64
	ns = int64(hour*3600+minute*60+second) * Second2Ns
	d := time.Duration(ns)
	return t.Add(d)
}

//获取子Map,需要做类型转换
func GetSubMap(key string, m map[string]interface{}) map[string]interface{} {
	sm, ok := m[key]
	if ok {
		return sm.(map[string]interface{})
	}
	return nil
}

//Map转 Struct ,需要通过中转
func Map2Struct(m map[string]interface{}, s interface{}) {
	bytes, err := json.Marshal(m)
	exception.CheckError(err)

	err = json.Unmarshal(bytes, &s)
	exception.CheckError(err)
}

// Struct 转 JSON, 与 json.Marshal() 的区别是将 sql.Null... 转换为 基本的 int,string,bool,float64 类型
// 在正则匹配中，必须加上 ? ,表示非贪婪模式，即最小匹配，否则最大匹配，就会匹配出最长的字符串
// var regNullString = regexp.MustCompile(`{"String":"(.*?)","Valid":(false|true)}`)
// var regNullN64 = regexp.MustCompile(`{("Float64"|"Int64"):([\d.]+?),"Valid":(false|true)}`)
// 统一一个语句完成
// 这个函数几乎可以不用了，直接对于存在 null 值的变量，使用指针解决，将会返回空指针，JSON 后会为 key : null ,前端可以处理
var regNull = regexp.MustCompile(`{("String"|"Float64"|"Int64"):(["]{0,1})(.*?)(["]{0,1}),"Valid":(false|true)}`)

func StructNullToJSON(s interface{}) string {
	data, err := json.Marshal(s)
	exception.CheckError(err)
	return regNull.ReplaceAllString(string(data), `$2$3$4`)
	//	ss := regNullString.ReplaceAllString(string(bs), `"$1"`)	//$1 表示 () 的位置
	//	return regNullN64.ReplaceAllString(ss, `$2`)
}

//将文本中的数字提取出来,主要用来解析怕从页面的数字解析
var regexInt = regexp.MustCompile(`([\d]+)`)

func Str2Int(s string) int {
	ints := regexInt.FindAllString(s, -1)
	if ints == nil {
		return 0
	}
	i, _ := strconv.Atoi(strings.Join(ints, ""))
	return i
}

//将只包含数字，逗号和小数点的字符串转为 Float64
func Str2Float64(s string) float64 {
	floats := regexFloat.FindAllString(s, -1)
	f, err := strconv.ParseFloat(strings.Join(floats, ""), 64)
	if err != nil {
		return 0
	}
	return f
}

// 从字符串中提取浮点数，其中字符串中含有币别符号
// 不同国家的小数点处理方式(使用欧元的，小数点是逗号表示，分节符反而用点号)
// DE IT ES FR	形如 EUR 32,99 或者 1.234,88
//此处不加逗号，目的是通过逗号分开数字
var regexFloat = regexp.MustCompile(`([\d.]+)`)

func Currency2Float64(s string) float64 {
	floats := regexFloat.FindAllString(s, -1)
	if floats == nil {
		return 0
	}
	if l := len(floats); l > 1 { //说明数字只间存在逗号（如果是数字间存在其他符号，非法的）
		// 1.234,88 返回  [1.234 88],所以只需要判断后边是两位即可
		// 不过如果严格来说，有多位小数的，不能采用这种方式，必须分国家处理了
		if len(floats[l-1]) == 2 { //欧元体系的，数字, 如果后边是3位的，为日本，日本无小数
			s = strings.Replace(strings.Join(floats[:l-1], ""), ".", "", -1) + "." + floats[l-1]
		} else {
			s = strings.Join(floats, "")
		}
	} else {
		s = floats[0]
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// 指针赋值语句，直接赋值因无内存地址，需要通过地址方式赋值
// *s.v = 100 或者 *v = 100 是不行的，因 v 此时是 nil
func PtrInt(v int) *int {
	return &v
}
func PtrInt64(v int64) *int64 {
	return &v
}
func PtrFloat64(v float64) *float64 {
	return &v
}
func PtrString(v string) *string {
	return &v
}
