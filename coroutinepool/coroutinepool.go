package coroutinepool

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Pool 创建协程池结构体
type Pool struct {
	state            int32          //协程池状态
	capacity         int32          //协程池容量 开启 worker 数量的上限
	running          int32          //协程池中 当前正在执行任务的 worker 数量
	lock             sync.Locker    //非阻塞锁--用于协程池访问 worker 数组
	workers          *workerQueue   //循环队列存储可用的worker
	workerPool       sync.Pool      //goWorker 池
	waiting          int32          //阻塞的goroutine 数量
	cond             *sync.Cond     //条件通知
	MaxBlockingTasks int            //当前最大阻塞的任务
	Wg               sync.WaitGroup //
}

var workerChanCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}
	return 1
}()

// NewPool 协程池构造函数 [size 协程池容量大小] [limitTask 限制任务数]
func NewPool(size int, limitTask int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrInvalidParams
	}
	//创建协程池实例
	p := &Pool{
		capacity:         int32(size),
		lock:             NewNonBlockingLock(), //非阻塞锁
		MaxBlockingTasks: limitTask,
	}
	//初始化workers 队列
	p.workers = newWorkerQueue(size)

	//初始化sync.pool
	p.workerPool.New = func() any {
		return &goWorker{
			pool: p,
			task: make(chan func(), workerChanCap),
		}
	}
	//初始化cond
	p.cond = sync.NewCond(p.lock)

	return p, nil
}

//------public method

// Submit pool 的提交任务方法
func (p *Pool) Submit(task func()) error {
	//异常处理
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var gw *goWorker
	if gw = p.obtainGoWorker(); gw == nil {
		return errObtainGoWorkerFailed
	}

	//任务放入goWorker 的任务通道
	gw.task <- task
	return nil
}

// Running 原子操作 获取正在运行的goroutine个数
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Waiting 获取当前等待运行的task
func (p *Pool) Waiting() int {
	return int(atomic.LoadInt32(&p.waiting))
}

// Cap 原子操作 获取当前pool的容量
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// IsClosed 判断当前协程池是否关闭
func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == CLOSED
}

// CloseAndRelease 关闭协程池 并且清空workers 队列
func (p *Pool) CloseAndRelease() {
	if !atomic.CompareAndSwapInt32(&p.state, OPENED, CLOSED) {
		//关闭失败
		return
	}
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()

	p.cond.Broadcast()
}

//-----private method 私有方法
//协程池添加正在运行的goroutine 个数
func (p *Pool) addRunning(num int) {
	atomic.AddInt32(&p.running, int32(num))
}

//增加等待task 数量
func (p *Pool) addWaiting(num int) {
	atomic.AddInt32(&p.waiting, int32(num))
}

//获取一个可用的goWorker 去执行任务
func (p *Pool) obtainGoWorker() (gw *goWorker) {
	//加锁访问 worker 队列
	p.lock.Lock()
	gw = p.workers.pop() //优先从队列中获取
	if gw != nil {
		p.lock.Unlock()
	} else if p.capacity == 0 || p.capacity > p.running {
		//如果队列为空 且协程池有多余的容量 则新建一个新的goWorker
		p.lock.Unlock()
		gw = p.workerPool.Get().(*goWorker)
		gw.run()
	} else {
	retry:
		//保持获取锁  等待有可执行任务的 goWorker
		if p.Waiting() >= p.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}

		p.addWaiting(1)
		p.cond.Wait() //等待有空闲goWorker 通知
		p.addWaiting(-1)

		//如果被通知是协程池已经关闭 则抛弃所有任务 退出
		if p.IsClosed() {
			p.lock.Unlock()
			return
		}

		//从队列中取不出goWorker
		if gw = p.workers.pop(); gw == nil {
			if p.Running() < p.Cap() { //正在运行的goWorker 小于 pool容量
				p.lock.Unlock()
				gw = p.workerPool.Get().(*goWorker)
				gw.run()
				return
			}
			goto retry //继续尝试获取goWorker
		}
		p.lock.Unlock()
	}
	return
}

//将 goWorker 返回pool 复用
func (p *Pool) revertWorker(gw *goWorker) bool {
	//如果正在运行的goroutine 大于 pool 的容量 或者 pool 已经关闭
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || p.IsClosed() {
		p.cond.Broadcast() //通知所有等待的中的task 放弃一轮的等待
		return false
	}
	p.lock.Lock()

	//双重检测 防止内存泄漏
	if p.IsClosed() {
		p.lock.Unlock()
		return false
	}

	err := p.workers.put(gw) //goWorker 放入 workerQueue 队列失败
	if err != nil {
		p.lock.Unlock()
		return false
	}
	//放入成功 则通知
	p.cond.Signal()
	p.lock.Unlock()
	return true
}
