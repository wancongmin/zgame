package alg

/*
	区间描述
*/

// IntervalTree 用于描述一个区间
type IntervalTree struct {
}

// interval 用于描述一个区间
type interval struct {
	min, max int64
	v        interface{}
}

// NewIntervalTree 用于新建立一个区间树
func NewIntervalTree() *IntervalTree {

}
