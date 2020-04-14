package mycsv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//根据Struct 的 csv 属性获取对应的字段的位置
//keys 表头数组
func GetStructFieldMap(keys []string, v interface{}) map[int]int {
	km := map[string]int{} //表头列的位置
	for i, kv := range keys {
		km[kv] = i
	}
	m := map[int]int{}
	//字段反射
	t := reflect.TypeOf(v).Elem()
	vt := reflect.ValueOf(v).Elem()
	ilen := t.NumField()
	for i := 0; i < ilen; i++ {
		vf := vt.Field(i)
		if vf.CanSet() == false {
			continue
		}
		csv := t.Field(i).Tag.Get("csv")
		if csv != "" {
			ki, ok := km[csv]
			if ok == true {
				m[i] = ki //Struct的字段列序与表头列序的对应关系
			}
		}
	}
	return m
}

//将数据转化为结构,只支持简单的数据类型，不支持 指针类型和 slice,对于这种，先对应后再处理
//keys 标题数组
//data 行值
//v 结构体
func SetStructValue(keyMap map[int]int, data []string, v interface{}) {
	t := reflect.ValueOf(v).Elem()
	for si, di := range keyMap {
		sd := data[di]
		fv := t.Field(si)
		switch fv.Kind() {
		case reflect.Int, reflect.Int64:
			if sd == "" {
				fv.SetInt(0)
			} else {
				id, err := strconv.ParseInt(sd, 10, 64)
				if err != nil {
					id = 0
				}
				fv.SetInt(id)
			}
			break
		case reflect.Float64, reflect.Float32:
			if sd == "" {
				fv.SetFloat(0)
			} else {
				fd, err := strconv.ParseFloat(sd, 64)
				if err != nil {
					return
				}
				fv.SetFloat(fd)
			}
			break
		default:
			fv.SetString(sd)
		}
	}
}

//根据数据和头的字段生成Map数据
func GetMap(keys []string, data []string, m map[string]string) {
	for i, key := range keys {
		m[key] = data[i]
	}
}

//生成表头映射关系
//将表头的空格换成下划线，全部小写
func GetKeyMap(header []string) {
	keys := "{"
	for _, key := range header {
		keys = keys + "\"" + strings.Replace(strings.ToLower(key), " ", "_", -1) + "\","
	}
	keys += "}"
	fmt.Println(keys)
}
