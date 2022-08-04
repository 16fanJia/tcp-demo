package coroutinepool

type workerQueue struct {
	items  []*goWorker //worker队列
	front  int         //头
	rear   int         //尾
	size   int         //大小
	isFull bool        //判断队列是否已满
}

//worker队列构造函数
func newWorkerQueue(size int) *workerQueue {
	return &workerQueue{
		items: make([]*goWorker, size),
		size:  size,
	}
}

//获取此时 可用 worker 循环队列的长度
func (w *workerQueue) len() int {
	if w.size == 0 {
		return 0
	}

	if w.front == w.rear {
		//队满
		if w.isFull {
			return w.size
		}
		//队空
		return 0
	}

	return (w.rear - w.front + w.size) % w.size
}

//判断队列是否为空
func (w *workerQueue) isEmpty() bool {
	return w.front == w.rear && !w.isFull
}

//将worker 放入队列
func (w *workerQueue) put(worker *goWorker) error {
	//队列异常处理 抛出异常
	if w.size == 0 {
		return errQueueLengthIsZero
	}
	if w.isFull {
		return errQueueIsFull
	}
	//worker 放入队列
	w.items[w.rear] = worker
	w.rear = (w.rear + 1) % w.size //队尾++

	//队满
	if (w.rear)%w.size == w.front {
		w.isFull = true
	}

	return nil
}

//队列弹出worker
func (w *workerQueue) pop() *goWorker {
	//队列空 弹出失败
	if w.isEmpty() {
		return nil
	}

	worker := w.items[w.front]
	w.front = (w.front + 1) % w.size

	//出队后 队列有空闲时间
	w.isFull = false

	return worker
}

//重置队列
func (w *workerQueue) reset() {
	if w.isEmpty() {
		return
	}
	//Releasing:
	//	if goW := w.pop(); goW != nil {
	//		//队列中还有worker 任务channel置空 然后被放入sync.pool
	//		goW.task <- nil
	//		goto Releasing
	//	}
	for {
		//for 循环一直从worker 队列中取worker
		goW := w.pop()
		if goW != nil {
			//队列中还有worker 任务channel置空 然后被放入sync.pool
			goW.task <- nil
		} else {
			break
		}
	}

	w.items = w.items[:0]
	w.size = 0
	w.front = 0
	w.rear = 0
}
