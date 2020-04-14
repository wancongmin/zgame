package config

import (
	"zinx/lib/exception"
	"encoding/json"
	"io/ioutil"
	"log"
)

var configMap map[string]interface{}

//自动执行配置文件加载
// init main 全局变量加载顺序；按引用关系，优先加载被引用的，
// 在同一个引用 package, 变量优先，init其次，main 最后
func init() {
	InitConfig()
}

/**
初始话Config,采用JSON存储配置信息，文件名必须为 config.json
*/
func InitConfig() {
	if configMap != nil {
		return
	}
	configMap = make(map[string]interface{})
	data, err := ioutil.ReadFile("config.json")
	exception.CheckError(err)

	err = json.Unmarshal(data, &configMap)
	exception.CheckError(err)
	log.Println("配置文件加载成功")
}

/**
根据Key值获取配置信息
*/
func GetConfig(key string, def string) string {
	v, ok := configMap[key]
	if ok {
		return v.(string)
	}
	return def
}

func GetIntConfig(key string, def int) int {
	v, ok := configMap[key]
	if ok {
		return int(v.(float64))
	}
	return def
}

/**
获取配置信息为Map的配置项
*/
func GetMapConfig(key string) map[string]interface{} {
	v, ok := configMap[key]
	if !ok {
		return nil
	}
	vm, ok := v.(map[string]interface{})
	if ok {
		return vm
	}
	return nil
}

//获取 Interface 配置，可以自己根据类型转换
func GetInterface(key string) interface{} {
	v, ok := configMap[key]
	if !ok {
		return nil
	}
	return v
}

//获取字符串数组
func GetConfigStrings(key string) (ret []string) {
	v, ok := configMap[key]
	if !ok {
		return ret
	}
	vi := v.([]interface{}) //该处不能直接 v.([]string)
	for _, s := range vi {
		ret = append(ret, s.(string))
	}
	return ret
}
