package main

import "sync"

/* sync.RWMutex (TODO: 源码) */

/*
	读写互斥锁 sync.RWMutex 是细粒度的互斥锁, 不限制资源的并发读, 但是读写,
	写写操作无法并行执行.

	常见服务的资源读写比例会非常高, 因为大多数读请求之间不会相互影响, 所有可以
	分离读写操作, 以此提高服务的性能

	type RWMutex struct {
		w           Mutex
		writerSem   uint32
		readerSem   uint32
		readerCount int32
		readerWait  int32
	}
	- w 复用互斥锁提供的能力;
	- writerSem 和 readerSem 分别用于写等待读和读等待写;
	- readerCount 存储了当前正在执行的读操作数量;
	- readerWait 表示当写操作被阻塞时等待的读操作个数;

	TODO:
	写操作使用 sync.RWMutex.Lock() 和 sync.RWMutex.Unlock()
	读操作使用 sync.RWMutex.RLock()  和 sync.RWMutex.RUnlock()


	* 写操作
	// Lock locks rw for writing.
	// If the lock is already locked for reading or writing,
	// Lock blocks until the lock is available.
	func (rw *RWMutex) Lock() {
		if race.Enabled {
			_ = rw.w.state
			race.Disable()
		}
		// First, resolve competition with other writers.
		rw.w.Lock()
		// Announce to readers there is a pending writer.
		r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
		// Wait for active readers.
		if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
			runtime_SemacquireMutex(&rw.writerSem, false, 0)
		}
		if race.Enabled {
			race.Enable()
			race.Acquire(unsafe.Pointer(&rw.readerSem))
			race.Acquire(unsafe.Pointer(&rw.writerSem))
		}
	}
	- 使用 RWMutex 结构体中的 sync.Mutex.Lock 阻塞后续的写操作; 因为互斥锁
			已经被获取, 所以其他的goroutine 在获取写锁时会进入自旋或者休眠;
	- 调用 sync/atomic.AddInt32 函数阻塞后续的读操作;
	- 如果仍有其他的 goroutine 持有互斥锁的读锁, 该 goroutine 会调用
			runtime.sync_runtime_SemacquireMutex 进入休眠状态, 等待所有读锁
			所有者执行结束后释放 writerSem 信号量唤醒当前 goroutine;

	写锁的释放会调用 sync.RWMutex.Unlock:
	// unlock unlocks rw for writing. it is a run-time error if rw is
	// not locked for writing on entry to unlock.
	//
	// as with mutexes, a locked rwmutex is not associated with a particular
	// goroutine. one goroutine may rlock (lock) a rwmutex and then
	// arrange for another goroutine to runlock (unlock) it.
	func (rw *rwmutex) unlock() {
		if race.enabled {
			_ = rw.w.state
			race.release(unsafe.pointer(&rw.readersem))
			race.disable()
		}

		// announce to readers there is no active writer.
		r := atomic.addint32(&rw.readercount, rwmutexmaxreaders)
		if r >= rwmutexmaxreaders {
			race.enable()
			throw("sync: unlock of unlocked rwmutex")
		}
		// unblock blocked readers, if any.
		for i := 0; i < int(r); i++ {
			runtime_semrelease(&rw.readersem, false, 0)
		}
		// allow other writers to proceed.
		rw.w.unlock()
		if race.enabled {
			race.enable()
		}
	}
	与加锁的过程相反, 解锁的过程如下:
	- 调用 sync/atomic.AddInt32 函数将 readercount 变回正数, 释放读锁;
	- 通过 for 循环释放所有因获取读而陷入等待的 goroutine;
	- 调用 sync.Mutex.Unlock 释放写锁

	* 获取写锁时会先阻塞写锁的获取, 后阻塞读锁的获取, 这种策略能够保证读操作
		不会因连续的写操作"饿死";



	* 读操作
	// RLock locks rw for reading.
	//
	// It should not be used for recursive read locking; a blocked Lock
	// call excludes new readers from acquiring the lock. See the
	// documentation on the RWMutex type.
	func (rw *RWMutex) RLock() {
		if race.Enabled {
			_ = rw.w.state
			race.Disable()
		}
		if atomic.AddInt32(&rw.readerCount, 1) < 0 {
			// A writer is pending, wait for it.
			runtime_SemacquireMutex(&rw.readerSem, false, 0)
		}
		if race.Enabled {
			race.Enable()
			race.Acquire(unsafe.Pointer(&rw.readerSem))
		}
	}
	通过 sync/atomic.AddInt32 将 readerCount 加1;
	- 如果返回负数, 其他 goroutine 获得了写锁, 当前 goroutine 就会调用
			runtime.sync_runtime_SemacquireMutex  陷入休眠等待锁的释放:
	 - 如果返回非负数, 则没有 goroutine 获得写锁, 当前方法成功返回;


	 调用 sybnc.RWMutex.RUnlock 释放读锁:
	// RUnlock undoes a single RLock call;
	// it does not affect other simultaneous readers.
	// It is a run-time error if rw is not locked for reading
	// on entry to RUnlock.
	func (rw *RWMutex) RUnlock() {
		if race.Enabled {
			_ = rw.w.state
			race.ReleaseMerge(unsafe.Pointer(&rw.writerSem))
			race.Disable()
		}
		if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
			// Outlined slow-path to allow the fast-path to be inlined
			rw.rUnlockSlow(r)
		}
		if race.Enabled {
			race.Enable()
		}
	}
	使用 sync/atomic.AddInt32 先减少正在读资源的 readerCount, 如:
	- 返回值大于等于0, 读锁直接解锁成功;
	- 如果返回值小于0, 有一个写操作正在进行, 这时使用 sync.RWMutex.rUnlockSlow
		方法减少获取锁的写操作等待的读操作数 readerWait, 并在所有读操作被释放
		后触发写操作的信号量 writerSem, 该信号量被触发时, 调度器就会唤醒尝试获
		取写锁的 goroutine;

	* 读写互斥锁在互斥锁上提供了额外的更细粒度的控制, 能够在读操作远远多于
		写操作时提升性能; (TODO: 实际的应用场景/练习实验)
*/
