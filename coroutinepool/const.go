package coroutinepool

import (
	"math"
)

const (
	DefaultPoolSize = math.MaxInt32 //默认协程池大小
)

//协程池状态
const (
	OPENED = iota //开启的
	CLOSED        //关闭的
)
