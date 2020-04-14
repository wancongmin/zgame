package util

import (
	"sync"
	"time"
)

/**
  系统自带的随机数,在高并发的情况下会产生延迟,因此,需要改善一下

*/
// Uint32 多协程调用安全,返回一个随机数
func Uint32() uint32 {
	v := rngPool.Get()
	if v == nil {
		v = &RNG{}
	}
	r := v.(*RNG)
	x := r.Uint32()
	rngPool.Put(r)
	return x
}

var rngPool sync.Pool

// Uint32n 返回一个随机数,返回一个在[0,maxN)之间的随机数
func Uint32n(maxN uint32) uint32 {
	x := Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

// RNG 多协程调用时不安全的
type RNG struct {
	x uint32
}

// Uint32 多协程调用不安全
func (r *RNG) Uint32() uint32 {
	for r.x == 0 {
		r.x = getRandomUint32()
	}
	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x
	return x
}

// Uint32n 返回一个在[0,maxN)之间的随机数
func (r *RNG) Uint32n(maxN uint32) uint32 {
	x := r.Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

// getRandomUint32 用于获取一个随机数
func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}
