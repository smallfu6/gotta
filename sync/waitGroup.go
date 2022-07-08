package main

/*
	sync.WaitGroup(TODO: 源码)
	sync.WaitGroup 可以等待一组 goroutine 返回; 一个常见的场景是批量发出
	RPC 或者 HTTP 请求
*/

import (
	"sync"
)

type Request struct{}

func main() {
	requests := []*Request{}
	wg := &sync.WaitGroup{}
	wg.Add(len(requests))

	for _, req := range requests {
		go func(r *Request) {
			defer wg.Done()
			// ...
		}(req)
	}

	wg.Wait()
}

/*
  TODO: 深入理解
  可以通过 sync.WaitGroup 将原本顺序执行的代码在多个 goroutine 中并发执行,
  加快程序处理的速度
										   --------------> Goroutine(Done)
                                         /
	WaitGroup  ---->    Goroutine(wait)  ----------------> Goroutine
	                                     \
										  ---------------> Goroutine

	// A WaitGroup must not be copied after first use.
	type WaitGroup struct {
		noCopy noCopy

		// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
		// 64-bit atomic operations require 64-bit alignment, but 32-bit
		// compilers do not ensure it. So we allocate 12 bytes and then use
		// the aligned 8 bytes in them as state, and the other 4 as storage
		// for the sema.
		state1 [3]uint32
	}
	* noCopy 保证 sync.WaitGroup 不会被开发者通过再赋值的方式复制;
		noCopy 是一个特殊的私有结构体,
		src/cmd/vendor/golang.org/x/tools/go/analysis/passes/copylock 包中的分析器
		会在编译期间检查被复制的变量中是否包含 sync.noCopy 或者实现了 Lock 和
		Unlock 方法, 如果包含该结构体或者实现了对应方法, 抛出错误.


		// func main() {
		// 	wg := sync.WaitGroup{}
		// 	yawg := wg
		// 	fmt.Println(wg, yawg)
		// }
		// 因为变量赋值或者调用函数时发生值复制导致分析器报错.
		// #command-line-arguments
		// ./waitGroup.go:56:10: assignment copies lock value to yawg: sync.WaitGroup contains sync.noCopy
		// ./waitGroup.go:57:14: call of fmt.Println copies lock value: sync.WaitGroup contains sync.noCopy
		// ./waitGroup.go:57:18: call of fmt.Println copies lock value: sync.WaitGroup contains sync.noCopy

	* state1 存储状态和信号量
		占12字节的数组, 该数组会存储当前结构体的状态, 在64位和32位机器上有所差异

		64bits         waiter       counter      sema
		32bits         sema         waiter       counter
		私有方法 sync.WaitGroup.state 能够从 state1 字段中取出它的状态和信号量
*/

/*
	sync.WaitGroup 对外提供的方法有: sync.WaitGroup.Add, sync.WaitGroup.Wait,
		sync.WaitGroup.Done;

	* ----------------------Add
	func (wg *WaitGroup) Add(delta int) {
		statep, semap := wg.state()
		if race.Enabled {
			_ = *statep // trigger nil deref early
			if delta < 0 {
				// Synchronize decrements with Wait.
				race.ReleaseMerge(unsafe.Pointer(wg))
			}
			race.Disable()
			defer race.Enable()
		}
		state := atomic.AddUint64(statep, uint64(delta)<<32)
		v := int32(state >> 32)
		w := uint32(state)
		if race.Enabled && delta > 0 && v == int32(delta) {
			// The first increment must be synchronized with Wait.
			// Need to model this as a read, because there can be
			// several concurrent wg.counter transitions from 0.
			race.Read(unsafe.Pointer(semap))
		}
		if v < 0 {
			panic("sync: negative WaitGroup counter")
		}
		if w != 0 && delta > 0 && v == int32(delta) {
			panic("sync: WaitGroup misuse: Add called concurrently with Wait")
		}
		if v > 0 || w == 0 {
			return
		}
		// This goroutine has set counter to 0 when waiters > 0.
		// Now there can't be concurrent mutations of state:
		// - Adds must not happen concurrently with Wait,
		// - Wait does not increment waiters if it sees counter == 0.
		// Still do a cheap sanity check to detect WaitGroup misuse.
		if *statep != state {
			panic("sync: WaitGroup misuse: Add called concurrently with Wait")
		}
		// Reset waiters count to 0.
		*statep = 0
		for ; w != 0; w-- {
			runtime_Semrelease(semap, false, 0)
		}
	}


	* ----------------------Wait
	// Wait blocks until the WaitGroup counter is zero.
	func (wg *WaitGroup) Wait() {
		statep, semap := wg.state()
		if race.Enabled {
			_ = *statep // trigger nil deref early
			race.Disable()
		}
		for {
			state := atomic.LoadUint64(statep)
			v := int32(state >> 32)
			w := uint32(state)
			if v == 0 {
				// Counter is 0, no need to wait.
				if race.Enabled {
					race.Enable()
					race.Acquire(unsafe.Pointer(wg))
				}
				return
			}
			// Increment waiters count.
			if atomic.CompareAndSwapUint64(statep, state, state+1) {
				if race.Enabled && w == 0 {
					// Wait must be synchronized with the first Add.
					// Need to model this is as a write to race with the read in Add.
					// As a consequence, can do the write only for the first waiter,
					// otherwise concurrent Waits will race with each other.
					race.Write(unsafe.Pointer(semap))
				}
				runtime_Semacquire(semap)
				if *statep != 0 {
					panic("sync: WaitGroup is reused before previous Wait has returned")
				}
				if race.Enabled {
					race.Enable()
					race.Acquire(unsafe.Pointer(wg))
				}
				return
			}
		}
	}

	* 总结:
	- sync.WaitGroup 必须在 sync.WaitGroup.Wait 方法返回之后才能重新使用;
	- Add 方法传入任意负数(需要保证计数器非负), 快速将计数器归零以唤醒等待的
		goroutine;
	- 可以同时有多个 goroutine 等待当前 sync.WaitGroup 计数器归零, 这些
		goroutine 会被同时唤醒;
*/
