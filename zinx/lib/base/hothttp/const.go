package hothttp

const (
	routerBucket = 64 // 路由桶的大小
)

// InterfaceHandle 用于描述一个对外的接口
type InterfaceHandle func(*Session)
