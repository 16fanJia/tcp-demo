package coroutinepool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	wg         sync.WaitGroup
	err        error
	PoolClient *Pool
)

func TestNewPool(t *testing.T) {
	PoolClient, err = NewPool(10, 100)
	if err != nil {
		return
	}
	defer PoolClient.CloseAndRelease()

	for i := 0; i < 1000; i++ {
		PoolClient.Wg.Add(1)
		if err = PoolClient.Submit(DemoFunc); err != nil {
			return
		}
	}
	PoolClient.Wg.Wait()
}

func DemoFunc() {
	defer PoolClient.Wg.Done()
	fmt.Println("hello world")

	fmt.Println("hello coroutinePool")
	time.Sleep(time.Second)
}
