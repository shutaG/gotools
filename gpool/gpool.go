package gpool

import (
	"sync"
)

// GPool 通过channel来限制同时并发的数量
// 1. 最小并发树限制为1
// 2. 可以通过add添加或者减少正在执行中的goroutine
// 3. 通过done来标识一个协程的执行结束
// 4. 通过wait等待所有协程执行结束
type GPool struct {
	queue chan int
	wg    *sync.WaitGroup
}

func NewGPool(size int) *GPool {
	if size <= 0 {
		size = 1
	}
	return &GPool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add 超过chan的容量，则会阻塞
func (p *GPool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

func (p *GPool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *GPool) Wait() {
	p.wg.Wait()
}
