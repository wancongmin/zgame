package encrypt

import (
	"encoding/base64"
	"strings"
)

//Base64 加密
func EncodeBase64(s string) string {
	data := []byte(s)
	s = base64.StdEncoding.EncodeToString(data)
	s = strings.Replace(s, "=", "", -1) //去掉 =
	return s[2:] + s[0:2]
}

//Base64解密
func DecodeBase64(s string) string {
	ln := len(s)
	s = s[ln-2:] + s[:ln-2]
	m := ln % 4
	if m > 0 {
		s += strings.Repeat("=", 4-m)
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(data)
}
