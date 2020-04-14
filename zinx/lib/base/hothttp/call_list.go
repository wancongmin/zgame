package hothttp

// execList 用于描述一个执行连
type execList struct {
	list []InterfaceHandle
}

// newExecList 用于一个执行链路
func newExecList() *execList {
	pInst := new(execList)
	return pInst
}
