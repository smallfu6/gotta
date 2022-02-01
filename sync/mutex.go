package main

/*
	sync.Mutex(TODO:源码)
*/

type Mutex struct {
	state int32  // 表示当前互斥锁的状态
	sema  uint32 //用于控制锁状态的信号量(TODO: 进程调度中的信号量)
}

// 只占8字节空间的结构体表示了 go 语言中的互斥锁.

/*
 *状态(TODO)


			waitersCount		starving		woken		locked
			  29bit              1bit            1bit        1bit

	默认情况下, 互斥锁的所有状态位都是0, int32 中的不同位分别表示不同状态:
	- mutexLocked: 互斥锁的锁定状态
	- mutexWoken: 从正常模式唤醒
	- mutexStarving: 当前互斥锁进入饥饿状态
	- waitersCount: 当前互斥锁上等待的 goroutine 个数

 * 正常模式和饥饿模式
	TODO: 根据 go1.9 的提交, 了解源码的实现, 并进行实验验证以下结论
	在正常模式下, 锁的等待者会按照先进先出的顺序获取锁, 但是刚被唤醒的goroutine
	与新创建的 goroutine 竞争时, 大概率不会获取到锁(?); 为了减少这种情况的出现,
	一旦 goroutine 超时1ms没有获取到锁, 它就会将当前互斥锁切换到饥饿模式, 防止
	部分 goroutine 被 "饿死".
	饥饿模式是通过 go1.9 通过提交 sync:make Mutex more fair 引入的优化(TODO),
	目的是保证互斥锁的公平性; 在饥饿模式下, 互斥锁会直接交给等待队列最前面的
	goroutine(快速把互斥锁资源提供给当前等待的 队列), 新的 goroutine 在该状态
	下不能获取锁, 也不会进入自旋状态, 只会在队列末尾等待, 如果一个
	goroutine 获取了互斥锁并且在队列末尾或者它等待的时间少于1ms, 当前的互斥
	锁就会切换会正常模式.

	相比而言, 正常模式下的互斥锁能够提供更好的性能, 而饥饿模式能避免 goroutine
	由于陷入等待无法获取锁而造成的高尾延时(TODO).


	TODO: 深入理解, 设计原理和目的
	自旋是一种多线程同步机制, 当前进程在进入自旋的过程中会一直保持cpu占用, 持续
	检查某个条件是否为真; 在多核 cpu 上, 自旋可以避免 goroutine 的切换, 使用
	恰当会对性能带来很大的增益, 但使用不当就会拖慢整个程序, 所以 goroutine 进入
	自旋的条件非常苛刻:
	- 互斥所只有在正常模式下才能进入自旋
	- runtime.sync_runtime_canSpin 需要返回 true:
		- 在有多个 cpu 的机器上运行
		- 当前 goroutine 为了获取该锁进入自旋的次数少于4
		- 当前机器至少存在一个正在运行的处理器 P 并且处理的运行队列为空
	一旦当前的 goroutine 能够进入自旋, 就会调用 runtime.sync_runtime_doSpin 和
	runtime.procyield 并执行 30 次 PAUSE 指令, 该指令只会占用 CPU 并消耗 CPU
	时间.
	func sync_runtime_doSpin() {
		procyield(active_spin_cnt)
	}

	// 汇编
	TEXT runtime.procyield(SB), NOSPLIT, $0-0
		MOVL           cycles+0(FP), AX
	again:
		PAUSE
		SUBL    $1, AX
		JNZ     again
		RET
*/

/*
 * 加锁和解锁
    使用 sync.Mutex.Lock 方法进行加锁, 当锁的状态是0时, 将 mutexLocked 设置成1

	// Lock locks m.
	// If the lock is already in use, the calling goroutine
	// blocks until the mutex is available.
	func (m *Mutex) Lock() {
		// Fast path: grab unlocked mutex.
		if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
			if race.Enabled {
				race.Acquire(unsafe.Pointer(m))
			}
			return
		}
		// Slow path (outlined so that the fast path can be inlined)
		m.lockSlow()
	}

	如果互斥锁的状态不是0, 就会调用 sync.Mutex.lockSlow 尝试通过自旋(spinning)
	等方式等待锁的释放.(TODO: 源码)
	func (m *Mutex) lockSlow() {
		var waitStartTime int64
		starving := false
		awoke := false
		iter := 0
		old := m.state
		for {
			//-----------------判断当前 goroutine 能否进入自旋等待互斥锁的释放
			// Don't spin in starvation mode, ownership is handed off to waiters
			// so we won't be able to acquire the mutex anyway.
			if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
				// Active spinning makes sense.
				// Try to set mutexWoken flag to inform Unlock
				// to not wake other blocked goroutines.
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					awoke = true
				}
				runtime_doSpin()
				iter++
				old = m.state
				continue
			}
			//-------互斥锁会根据上下文计算当前互斥锁的最新状态---------------
			// 几个不同的条件会分别更新 state 字段中存储的不同信息, mutexLocked,
			// mutexStarving, mutexWoken, mutexWaiterShift

			new := old
			// Don't try to acquire starving mutex, new arriving goroutines must queue.
			if old&mutexStarving == 0 {
				new |= mutexLocked
			}
			if old&(mutexLocked|mutexStarving) != 0 {
				new += 1 << mutexWaiterShift
			}
			// The current goroutine switches mutex to starvation mode.
			// But if the mutex is currently unlocked, don't do the switch.
			// Unlock expects that starving mutex has waiters, which will not
			// be true in this case.
			if starving && old&mutexLocked != 0 {
				new |= mutexStarving
			}
			if awoke {
				// The goroutine has been woken from sleep,
				// so we need to reset the flag in either case.
				if new&mutexWoken == 0 {
					throw("sync: inconsistent mutex state")
				}
				new &^= mutexWoken
			}
			//------------计算了新的互斥锁状态后, 会使用 CAS 函数
			// CompareAndSwapInt32 更新状态 ----------------------------------
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				if old&(mutexLocked|mutexStarving) == 0 {
					break // locked the mutex with CAS
				}
				// If we were already waiting before, queue at the front of the queue.
				queueLifo := waitStartTime != 0
				if waitStartTime == 0 {
					waitStartTime = runtime_nanotime()
				}
				runtime_SemacquireMutex(&m.sema, queueLifo, 1)
				starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
				old = m.state
				if old&mutexStarving != 0 {
					// If this goroutine was woken and mutex is in starvation mode,
					// ownership was handed off to us but mutex is in somewhat
					// inconsistent state: mutexLocked is not set and we are still
					// accounted as waiter. Fix that.
					if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
						throw("sync: inconsistent mutex state")
					}
					delta := int32(mutexLocked - 1<<mutexWaiterShift)
					if !starving || old>>mutexWaiterShift == 1 {
						// Exit starvation mode.
						// Critical to do it here and consider wait time.
						// Starvation mode is so inefficient, that two goroutines
						// can go lock-step infinitely once they switch mutex
						// to starvation mode.
						delta -= mutexStarving
					}
					atomic.AddInt32(&m.state, delta)
					break
				}
				awoke = true
				iter = 0
			} else {
				old = m.state
			}
			// 如果没有通过 CAS 获得锁, 会调用runtime.sync_runtime_SemacquireMutex
			// 通过信号量保证资源不会被两个 goroutine 获取;
			// runtime.sync_runtime_SemacquireMutex 会在方法中不断尝试获取锁并陷
			// 入休眠等待信号量释放, 一旦当前 goroutine 可以获取信号量, 它就会
			// 立刻返回, sync.Mutex.Lock 的剩余代码也会继续执行.
		   //-----------------------------------------------------------------
		}

		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
	}

	- 在正常模式下, lockSlow 会设置唤醒和饥饿标记, 重置迭代次数并重新执行获取锁
		的循环;
	- 在饥饿模式下, 当前goroutine 会获得互斥锁, 如果等待队列中只存在当前
		goroutine, 互斥锁还会从饥饿模式中退出.

*/
