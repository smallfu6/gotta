package main

import "sync"

/*
	在某些场景, 只是希望一个任务有单一的执行者, 而不像计数器场景那样所有
	Goroutine都执行成功; 后来的Goroutine在抢锁失败后, 需要放弃其流程,
	这时候就需要尝试锁(try lock)了;

	尝试锁: TODO: 熟悉 sync.Mutex 的 tryLock 方法
	如果加锁成功执行后续流程, 如果加锁失败也不会阻塞, 而会直接返回加锁的结果;
	在Go语言中可以用大小为1的通道模拟尝试锁; 也可以使用标准库函数CAS实现相同
	的功能且成本更低;

	在单机系统中, 尝试锁并不是一个好选择, 因为大量的Goroutine抢锁可能会导致
	CPU无意义的资源浪费; 有一个专有名词用来描述这种抢锁的场景——活锁, 指的是
	程序看起来在正常执行, 但实际上CPU周期被浪费在抢锁而非执行任务上, 从而
	导致程序整体的执行效率低下; 活锁的问题定位起来要麻烦很多, 所以在单机场
	景下, 不建议使用这种锁; 本节中模拟实现简单的尝试锁;

	基于 Zookeeper 实现分布式锁 ./zookeeper.go
	基于 etcd 实现分布式锁 ./etcd.go
	TODO: 分布式系统,分布式锁

*/

// Lock 尝试锁
type Lock struct {
	c chan struct{}
}

// NewLock 生成一个尝试锁
func NewLock() Lock {
	var l Lock
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

// Lock 锁住尝试锁返回加锁结果
func (l Lock) Lock() bool {
	lockResult := false
	select {
	case <-l.c:
		lockResult = true
	default:
	}
	return lockResult
}

func (l Lock) Unlock() {
	l.c <- struct{}{}
}

var counter int

func main() {
	var l = NewLock()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Lock() {
				println("lock failed")
				return
			}
			counter++
			println("current counter", counter)
			l.Unlock()
			// 限定每个Goroutine只有成功执行了Lock才会继续执行后续逻辑,
			// 因此在Unlock时可以保证Lock结构体中的通道一定是空, 从而不会阻塞,
			// 也不会失败
		}()
	}
	wg.Wait()
}
