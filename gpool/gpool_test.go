package gpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_GPool_Add(t *testing.T) {
	pool := NewGPool(0)
	if cap(pool.queue) != 1 {
		t.Errorf("Expected capacity of 1, got %d", cap(pool.queue))
	}
}

func Test_GPool_Done(t *testing.T) {
	pool := NewGPool(2)
	pool.Add(2)
	pool.Done() // 完成一个协程
	if len(pool.queue) != 1 {
		t.Errorf("Expected 1 goroutines in the queue after done, got %d", len(pool.queue))
	}
}

func Test_GPool_Wait(t *testing.T) {
	pool := NewGPool(2)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pool.Add(1)
		pool.Done()
		wg.Done()
	}()
	wg.Wait() // 等待协程完成
	if len(pool.queue) != 0 {
		t.Errorf("Expected 0 goroutines in the queue after wait, got %d", len(pool.queue))
	}
}

func Test_GPool(t *testing.T) {
	pool := NewGPool(2) // 限制同时并发2个协程
	var currentCount int32
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pool.Add(1) // 尝试添加一个协程到池中
			fmt.Println(i, time.Now())
			newCount := atomic.AddInt32(&currentCount, 1)
			if newCount > 2 {
				t.Errorf("More than 2 goroutines are running concurrently")
			}
			time.Sleep(1 * time.Second)
			// 模拟一些工作
			pool.Done() // 标记此协程已完成
			atomic.AddInt32(&currentCount, -1)
		}(i)
	}

	wg.Wait() // 等待所有协程完成
	finalCount := atomic.LoadInt32(&currentCount)
	if finalCount != 0 {
		t.Errorf("Not all goroutines have finished execution")
	}
}
