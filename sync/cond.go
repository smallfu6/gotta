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


	sync.Cond 不是常用的同步机制, 但是在条件长时间无法满足时, 与使用 for{} 进行
	忙碌等待相比, sync.Cond 能够让出处理器的使用权, 提高 cpu 的利用率, 使用时
	需注意:
		- sync.Cond.Wait 在调用之前一定要先获取互斥锁, 否则会触发程序崩溃
	    - sync.Cond.Signal 唤醒的 goroutine 都是队列最前面, 等待最久的 goroutine
		- sync.Cond.Broadcast 会按照一定顺序广播通知等待的全部 goroutine

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

	在 sync.notifyList 结构体中, head 和 tail 分别指向链表的头和尾, wait 和
	notify 分别表示正在等待的和已经通知到的 goroutine 索引:
	type notifyList struct {
		wait   uint32
		notify uint32
		lock   uintptr // key field of the mutex
		head   unsafe.Pointer
		tail   unsafe.Pointer
	}
*/
