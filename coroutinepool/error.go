package coroutinepool

import "errors"

var (
	ErrInvalidParams        = errors.New("参数无效！")
	ErrPoolClosed           = errors.New("此协程池已经关闭！")
	errQueueLengthIsZero    = errors.New("队列长度为0")
	errQueueIsFull          = errors.New("队列已满,put 失败")
	errObtainGoWorkerFailed = errors.New("获取可用工作协程失败")
)
