package main

import (
	"log"
	"sync"
	"time"
)

/*
	sync.Mutex(TODO:源码)
	// A Mutex is a mutual exclusion lock.
	// The zero value for a Mutex is an unlocked mutex.

	// A Mutex must not be copied after first use.
	type Mutex struct {
	   state int32  // 表示当前互斥锁的状态
	   sema  uint32 // 用于控制锁状态的信号量(TODO: 进程调度中的信号量)
	}
	sync包中类型的实例在首次使用后被复制得到的副本一旦再被使用将
	导致不可预期的结果, 为此在使用sync包中类型时, 推荐通过闭包方
	式或传递类型实例(或包裹该类型的类型实例)的地址或指针的方式使用
*/

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
			if old&mutexStarving == 0 { // TODO:  语法
				new |= mutexLocked   // TODO: 语法
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
				new &^= mutexWoken  // TODO: 语法
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




	相比加锁, 解锁过程 sync.Mutex.Unlock 较简单:
	func (m *Mutex) Unlock() {
		if race.Enabled {
			_ = m.state
			race.Release(unsafe.Pointer(m))
		}

		// Fast path: drop lock bit.
		new := atomic.AddInt32(&m.state, -mutexLocked)
		if new != 0 {
			// Outlined slow path to allow inlining the fast path.
			// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
			m.unlockSlow(new)
		}
	}
	该过程会先使用 atomic.AddInt32 函数快速解锁:
	- 如果该函数返回的新状态为0, 当前 goroutine 成功解锁;
	- 如果新状态不等于0, 会调用 unlockSlow 慢速解锁;

	unlockSlow 会先校验锁状态的合法性--如果当前互斥锁已经被解锁了, 会直接抛出
	异常 "sync: unlock of unlocked mutex" 中止当前程序.
	func (m *Mutex) unlockSlow(new int32) {
		if (new+mutexLocked)&mutexLocked == 0 {
			throw("sync: unlock of unlocked mutex")
		}
		if new&mutexStarving == 0 {
			old := new
			for {
				// If there are no waiters or a goroutine has already
				// been woken or grabbed the lock, no need to wake anyone.
				// In starvation mode ownership is directly handed off from unlocking
				// goroutine to the next waiter. We are not part of this chain,
				// since we did not observe mutexStarving when we unlocked the mutex above.
				// So get off the way.
				if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
					return
				}
				// Grab the right to wake someone.
				new = (old - 1<<mutexWaiterShift) | mutexWoken
				if atomic.CompareAndSwapInt32(&m.state, old, new) {
					runtime_Semrelease(&m.sema, false, 1)
					return
				}
				old = m.state
			}
		} else {
			// Starving mode: handoff mutex ownership to the next waiter, and yield
			// our time slice so that the next waiter can start to run immediately.
			// Note: mutexLocked is not set, the waiter will set it after wakeup.
			// But mutex is still considered locked if mutexStarving is set,
			// so new coming goroutines won't acquire it.
			runtime_Semrelease(&m.sema, true, 1)
		}
	}
	正常模式下:
	- 如果互斥锁不存在等待者, 或者互斥锁的 mutexLocked, mutexStarving, mutexWoken
		状态不都为 0, 当前方法可以直接返回, 不需要唤醒其他等待者;
	- 如果存在等待者, 会通过 sync.runtime_Semrelease 唤醒等待者并移交锁的所有权
	饥饿模式下:
	直接调用 sync.runtime_Semrelease, 将当前锁交给下一个正在尝试获取锁的等待者,
	等待者被唤醒后会得到锁, 这时互斥锁不会退出饥饿状态;

*/

type foo struct {
	n int
	sync.Mutex
}

func main() {
	f := foo{n: 17}
	go func(f foo) {
		for {
			log.Println("g2: try to lock foo...")
			f.Lock()
			log.Println("g2: lock foo ok")

			time.Sleep(3 * time.Second)

			f.Unlock()
			log.Println("g2: unlock foo ok")
		}
	}(f)

	f.Lock()
	log.Println("g1: lock foo ok")

	go func(f foo) {
		for {
			log.Println("g3: try to lock foo...")
			f.Lock()
			log.Println("g3: lock foo ok")

			time.Sleep(5 * time.Second)

			f.Unlock()
			log.Println("g3: unlock foo ok")
		}
	}(f)

	time.Sleep(1000 * time.Second)
	f.Unlock()
	log.Println("g1: unlock foo ok")
}

// 2022/07/06 12:13:26 g1: lock foo ok
// 2022/07/06 12:13:26 g3: try to lock foo...
// 2022/07/06 12:13:26 g2: try to lock foo...
// 2022/07/06 12:13:26 g2: lock foo ok
// 2022/07/06 12:13:29 g2: unlock foo ok
// 2022/07/06 12:13:29 g2: try to lock foo...
// 2022/07/06 12:13:29 g2: lock foo ok
// 2022/07/06 12:13:32 g2: unlock foo ok
// 2022/07/06 12:13:32 g2: try to lock foo...
// 2022/07/06 12:13:32 g2: lock foo ok
// 2022/07/06 12:13:35 g2: unlock foo ok
// 2022/07/06 12:13:35 g2: try to lock foo...
// 2022/07/06 12:13:35 g2: lock foo ok

/*
	创建了两个goroutine: g2和g3; 运行的结果显示: g3阻塞在加锁操作上了, 而按g2
	则按预期正常运行; g2和g3的差别就在于g2是在互斥锁首次使用之前创建的, 而g3
	则是在互斥锁执行完加锁操作并处于锁定状态之后创建的, 并且程序在创建g3的时候
	复制了foo的实例(包含sync.Mutex的实例)并在之后使用了这个副本;


	对Mutex实例的复制即是对两个整型字段的复制, 在初始状态下, Mutex实例处于
	Unlocked状态(state和sema均为0); g2复制了处于初始状态的Mutex实例,
	副本的state和sema均为0, 这与g2自定义一个新的Mutex实例无异, 这决定了g2
	后续可以按预期正常运行;

	后续主程序调用了Lock方法, Mutex实例变为Locked状态(state字段值
	为sync.mutex-Locked); 而此后g3创建时恰恰复制了处于Locked状态的
	Mutex实例(副本的state字段值亦为sync.mutexLocked); 因此g3
	再对其实例副本调用Lock方法将会导致其进入阻塞状态(也是死锁状态, 因为
	没有任何其他机会调用该副本的Unlock方法了, 并且Go不支持递归锁）
*/
