package coroutinepool

import (
	"runtime"
	"sync/atomic"
)

//非阻塞锁
type nonBlockingLock uint32

const maxBackoff = 16 //最大重试次数

// NewNonBlockingLock 构造函数 生成非阻塞锁实例
func NewNonBlockingLock() *nonBlockingLock {
	return new(nonBlockingLock)
}

func (nbl *nonBlockingLock) Lock() {
	//利用指数退避算法
	backoff := 1
	//利用CAS 尝试获取锁
	//取锁失败
	for !atomic.CompareAndSwapUint32((*uint32)(nbl), 0, 1) {
		//backoff = {1,2,4,8,16........,2^16-1}
		for i := 0; i < backoff; i++ {
			runtime.Gosched() //循环 让出cpu时间片
		}
		if backoff < maxBackoff {
			backoff <<= 1 //指数增长
		}
	}

}

func (nbl *nonBlockingLock) Unlock() {
	//解锁
	atomic.StoreUint32((*uint32)(nbl), 0)
}
