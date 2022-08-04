package coroutinepool

type goWorker struct {
	pool *Pool       //那个协程池拥有这个worker
	task chan func() //任务通道
}

//运行worker 开始一个goroutine 去运行任务
func (g *goWorker) run() {
	//新起一个协程 协程池正在运行的个数加一
	g.pool.addRunning(1)

	go func() {
		defer func() {
			//正在运行goroutine 的个数 -1
			g.pool.addRunning(-1)

			//没有任务的 worker 会被放回sync.pool中
			g.pool.workerPool.Put(g)
			g.pool.cond.Signal() //通知阻塞的task 可以获取goWorker来执行
		}()

		// 循环监听任务列表，一旦有任务立马取出任务运行
		for f := range g.task {
			if f == nil {
				return
			}
			f()
			//将运行完的goWorker 放入workers 队列
			if ok := g.pool.revertWorker(g); !ok {
				return
			}
		}
	}()

}
