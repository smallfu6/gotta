package main

/*
	channel 是 go 核心的数据结构和 goroutine 之间的通信方式, channel 是支撑 go
	语言高性能并发编程模型的重要结构;

	在很多主流编程语言中, 多个线程传递数据的方式一般是共享内存, 为了解决线程
	竞争, 需要限制同一时间能够读写这些变量的线程数量, 但这与 go 语言倡导的设计
	并不相同;

	虽然 go 也可以使用共享内存加互斥锁进行通信, 但是 go 提供了一种不同的并发
	模型--- 通信顺序进程(Communicating sequential processes, CSP); goroutine
	和 Channel 分别对应 CSP 中的实体和传递信息的媒介, goroutine 之间会通过
	channel 传递数据;


			Thread1 ---->  内存 ---->  Thread2
			多线程使用共享内存传递数据


			Goroutine ----->  Channel  -----> Goroutine
			goroutine 使用 channel 传递数据
		两个 goroutine 独立运行, 并不存在直接关联, 但是能通过 channel 间接完成
		通信;

		CSP 思想: "不要通过共享内存的方式进行通信, 而应该通过通信的方式共享内存"


	TODO:源码
	go 采用 csp 并发编程思想, 将通道作为协程之间交流的原语, 屏蔽了传统多线
	程编程中底层实现的诸多细节; 产生了很多经典的并发模型(ping-pong, fan-in,
	fan-out, pipeline);
	通道底层是用锁实现的环形队列, goroutine 在读取和写入时如果不能直接操作则
	会被放入等待队列中陷入休眠状态;  借助 go 运行时的调度器, 通道不会堵塞
	程序的执行(协程阻塞, 不就是堵塞了程序的执行?), 并且协程能够在需要时被唤醒;

	在实际中, 一个协程通常会处理多个通道, 我们不希望由于一个通道的读写而陷入
	阻塞, 影响其他通道的正常读写; 因此在实际中, 更多会使用 select 多路复用的
	机制同时监听多个通道是否准备就绪; select 在底层会锁住所有的通道并采取随
	机的方式保证公平的遍历所有通道;(TODO: select 源码)


	1.先进先出
		channel 的收发操作均遵循先进先出(FIFO)设计

	2.无锁 channel
		锁是常见的并发控制技术, 一般将锁分为"乐观锁" 和 "悲观锁", 即乐观并发
		控制和悲观并发控制; 无锁(lock-free) 队列是使用乐观并发控制的队列;
		注意, 乐观锁和悲观锁并不是真正的锁, 只是一种并发控制思想;
		乐观并发控制本质上是基于验证的协议, 使用原子指令 CAS(compare-and-swap,
		或 compare-and-set) 在多线程间同步数据, 无锁队列的实现也依赖这一原子
		指令; 因为性能原因, go 还未提供无锁 channel

	channel 在运行时的内部表示是 runtime.hchan, 该结构体中包含了用于保护
	成员变量的互斥锁; 可以说 channel 是一个用于同步和通信的有锁队列,
	使用互斥锁解决程序中可能存在的线程竞争问题很常见, 可以相对容易的实现
	有锁队列;


	数据结构
	go 的 channel 在运行时使用 runtime.hchan 结构体表示
	type hchan struct {
		qcount   uint           // total data in the queue
		dataqsiz uint           // size of the circular queue
		buf      unsafe.Pointer // points to an array of dataqsiz elements
		elemsize uint16
		closed   uint32
		elemtype *_type // element type
		sendx    uint   // send index
		recvx    uint   // receive index
		recvq    waitq  // list of recv waiters
		sendq    waitq  // list of send waiters

		// lock protects all fields in hchan, as well as several
		// fields in sudogs blocked on this channel.
		//
		// Do not change another G's status while holding this lock
		// (in particular, do not ready a G), as this can deadlock
		// with stack shrinking.
		lock mutex
	}

	sendq 和 recvq 存储了当前 channel 由于缓冲区空间不足而阻塞的 goroutine
	列表, 这些等待队列使用双向链表 runtime.waitq 表示, 链表中所有元素
	都是 runtime.sudog 结构;
	type waitq struct {
		first *sudog
		last  *sudog
	}

	对于有缓存的通道, 存储在 buf 中的数据虽然是线性的数组, 但是用数组和序号
	recvx, recvq 模拟了一个循环队列

*/

/*
	TODO: 源码
	通道的初始化在运行时调用了 makechan 函数
	func makechan(t *chantype, size int) *hchan {
		elem := t.elem

		// compiler checks this but be safe.
		if elem.size >= 1<<16 {
			throw("makechan: invalid channel element type")
		}
		if hchanSize%maxAlign != 0 || elem.align > maxAlign {
			throw("makechan: bad alignment")
		}

		mem, overflow := math.MulUintptr(elem.size, uintptr(size))
		if overflow || mem > maxAlloc-hchanSize || size < 0 {
			panic(plainError("makechan: size out of range"))
		}

		// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
		// buf points into the same allocation, elemtype is persistent.
		// SudoG's are referenced from their owning thread so they can't be collected.
		// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
		var c *hchan
		switch {
		case mem == 0:
			// Queue or element size is zero.
			c = (*hchan)(mallocgc(hchanSize, nil, true))
			// Race detector uses this location for synchronization.
			c.buf = c.raceaddr()
		case elem.ptrdata == 0:
			// Elements do not contain pointers.
			// Allocate hchan and buf in one call.
			c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
			c.buf = add(unsafe.Pointer(c), hchanSize)
		default:
			// Elements contain pointers.
			c = new(hchan)
			c.buf = mallocgc(mem, elem, true)
			// 当通道的元素中包含指针时, 需要单独分配内存空间, 因为当元素
			// 中包含指针时, 需要单独分配空间才能正常进行垃圾回收;
		}

		c.elemsize = uint16(elem.size)
		c.elemtype = elem
		c.dataqsiz = uint(size)
		lockInit(&c.lock, lockRankHchan)

		if debugChan {
			print("makechan: chan=", c, "; elemsize=", elem.size, "; dataqsiz=", size, "\n")
		}
		return c
	}

	(TODO: 阅读源码, 用 go 模拟实现一个简单的chan)
	通道写入原理
	发送元素时, 分成3种不同的情况: (以 c <- 5为例)
	- 有正在等待的读取协程
		hchan 的 recvq 字段存储了正在等待的协程链表, 每个协程对应一个 sudog
		结构, 是对协程的封装, 包含了准备获取的协程中的元素指针等; 当有读取的
		协程正在等待时, 直接从读取的协程链表中取第一个协程, 并将元素直接复制到
		对应的协程中, 再唤醒被阻塞的协程;(TODO: 唤醒原理)



				recvq ----> sudog1 ----> sudog2 ----> sudog3

 										|
										|
									   \|/

		recvq ----> sudog2 ----> sudog3     chan<-5 ------> send -------> sudog1

	- 缓冲区有空余
		如果队列中没有正在等待的协程, 但是该通道是带缓冲区的, 并且当前缓冲区
		没有满, 则向当前缓冲区中写入当前元素

	- 缓冲区无空余
		如果当前通道无缓冲区或者当前缓冲区已经满了, 则代表当前协程的 sudog 结构
		需要放入 sendq 链表末尾中, 并且当前协程陷入休眠状态, 等待被唤醒重新执行;

		sendq ----> sudog1 ----> sudog2 ----> sudog3(chan<-5)

	通道读取原理:
		- 有正在等待的写入协程
			当有正在等待的写入协程时, 直接从等待的写入协程链表中获取第一个
			协程(通道不带缓冲区的情况), 并将写入的元素直接复制到当前协程中,
			再唤醒被阻塞的写入协程;
		- 缓冲区有元素
			如果队列中没有正在等待的写入协程, 但是该通道是带缓冲区的, 并且当前
			缓冲区中有数据, 则读取该缓冲区中的数据, 并写入当前的读取协程中;
		- 缓冲区无元素
			如果当前通道无缓冲区或者当前缓冲区已经空了, 则代表当前协程的 sudog
			结构需要放入到 recvq 链表末尾, 并且当前协程陷入休眠状态, 等待被
			唤醒重新执行;

	select
	当 select 的分支都不满足时, 意味着当前协程需要进入休眠状态并等待 select 中
	至少有一个通道被唤醒; 不管是读取通道还是写入通道都需要创建一个新的 sudog
	并将其放入指定通道的等待队列, 之后当前协程将进入休眠状态;
	当 select 中的任意一个通道不再阻塞时, 当前协程将被唤醒; 要注意, 最后需要将
	sudog 结构体在其他通道的等待队列中出栈, 因为当前协程已经能够正常运行, 不再
	需要被其他通道唤醒;(TODO: 模拟简单实现)

	select 结构的执行过程与实现原理
	首先在编译期间, Go 语言会对 select 语句进行优化, 会根据 select 中 case 的
	不同选择不同的优化路径:
	1. 空的 select 语句会被转换成调用 runtime.block 直接挂起当前 Goroutine
	2. 如果 select 语句中只包含一个 case, 编译器会将其转换成 if ch == nil {
		block }; n; 表达式
		-. 首先判断操作的 Channel 是不是空的
		-. 然后执行 case 结构中的内容;
	3. 如果 select 语句中只包含两个 case 并且其中一个是 default, 那么会使用
		runtime.selectnbrecv 和 runtime.selectnbsend 非阻塞地执行收发操作;
	4. 在默认情况下会通过 runtime.selectgo 获取执行 case 的索引, 并通过多
		个 if 语句执行对应 case 中的代码

*/
