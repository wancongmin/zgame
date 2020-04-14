package util

// StringToInt 用于将字符串转换为Int值
func StringToInt(str string) int {
	length := len(str)
	sum := 0
	for i := 0; i < length; i++ {
		sum += int(str[i])
	}
	return sum
}
