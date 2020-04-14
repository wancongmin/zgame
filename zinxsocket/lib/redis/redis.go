package redis

import (
	"bangseller.com/lib/config"
	"bangseller.com/lib/exception"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var Redis *redis.Client

//常量定义
const (
	redisSection  = "Redis" //Redis配置节点名称
	redisAddr     = "Addr"
	redisPassword = "Password"
	redisDB       = "DB"
)

func InitRedis() {
	m := config.GetMapConfig(redisSection)
	if m == nil {
		return //为配置Redis信息
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     m[redisAddr].(string),
		Password: m[redisPassword].(string), // password set
		DB:       int(m[redisDB].(float64)), // use default DB
		PoolSize: 25,
	})

	pong, err := Redis.Ping().Result()
	exception.CheckError(err)
	log.Println("Redis 初始化成功,", pong)
}

//保存结构信息
func SetStruct(key string, s interface{}, expiration time.Duration) {
	data, err := json.Marshal(s)
	exception.CheckError(err)
	cmd := Redis.Set(key, data, expiration)
	exception.CheckError(cmd.Err())
}

//获取存取的结构信息
func GetStruct(key string, s interface{}) bool {
	data, err := Redis.Get(key).Bytes()
	if err == redis.Nil || len(data) == 0 {
		return false
	}
	exception.CheckError(err)
	err = json.Unmarshal(data, s)
	exception.CheckError(err)
	return true
}
