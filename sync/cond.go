package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

/*
	TODO: 源码
	条件变量 sync.Cond 可以让一组 goroutine 在满足特定条件时被唤醒,
    每个 sync.Cond 结构体在初始化时都需要传入一个互斥锁
	TODO: 应用场景
*/

var status int64

func main() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		// 10 个 goroutine 通过 sync.Cond.Wait 等待特定条件满足
		go listen(c)
	}

	time.Sleep(1 * time.Second)
	// 1个 goroutine 会调用 sync.Cond.Broadcast 唤醒所有陷入等待的 goroutine
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	c.L.Unlock()
}

func listen(c *sync.Cond) {
	c.L.Lock()
	for atomic.LoadInt64(&status) != 1 {
		c.Wait()
	}
	fmt.Println("listen")
	c.L.Unlock()
}

/*
	// A Cond must not be copied after first use.
	type Cond struct {
		noCopy noCopy

		// L is held while observing or changing the condition
		L Locker

		notify  notifyList
		checker copyChecker
	}
	noCopy: 保证结构体不会在编译期间复制
	L: 用于保护内部的 notify 字段, Locker 接口类型的变量
	notify: 一个 goroutine 的链表, 是实现同步机制的核心结构
	copyChecker: 用于禁止运行期间发生的复制

*/
