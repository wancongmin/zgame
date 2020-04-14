package request

import (
	"net/url"
	"strconv"
	"strings"
)

/**
获取指定key值为map，形如 key[..]=value1&key[...]=value2&...
*/
func FormMap(key string, values url.Values) map[string]string {
	m := make(map[string]string)

	key += "["
	kv := ""
	for k, v := range values {
		if strings.Index(k, key) == -1 {
			continue
		}
		if len(v) == 1 {
			kv = v[0]
		} else {
			kv = strings.TrimSpace(strings.Join(v, ""))
		}
		if kv == "" {
			continue
		}

		k = k[len(key) : len(k)-1]
		m[k] = kv
	}

	return m
}

/**
获取指定key值为map[int]map[string][string]，形如 key[0][..]=value01&key[0][...]=value02&key[1][..]=value11&key[1][...]=value12...
考虑到[0],[1],[2],..不一定连续，以及顺序也不定，以map返回
*/
func FormMap2(key string, uv url.Values) map[int]map[string]string {
	mapret := make(map[int]map[string]string)
	var m map[string]string
	var ok bool

	key += "["
	keyLen := len(key)
	kv := ""
	for k, v := range uv {
		if strings.Index(k, key) == -1 {
			continue
		}
		if len(v) == 1 {
			kv = v[0]
		} else {
			kv = strings.TrimSpace(strings.Join(v, ""))
		}
		if kv == "" {
			continue
		}
		keys := strings.Split(k, "][")
		if len(keys) < 2 {
			continue
		}

		ik, _ := strconv.Atoi(keys[0][keyLen:]) //序号
		m, ok = mapret[ik]
		if ok == false {
			m = make(map[string]string)
			mapret[ik] = m
		}

		k = keys[1][0 : len(keys[1])-1]
		m[k] = kv
	}

	return mapret
}
