package main

/*
	sync.Pool是一个数据对象缓存池, 有如下特点:
	- 是goroutine并发安全的, 可以被多个goroutine同时使用;
	- 放入该缓存池中的数据对象的生命是暂时的, 随时都可能被垃圾回收掉;
	- 缓存池中的数据对象是可以重复利用的, 可以在一定程度上降低数据对象
		重新分配的频度, 减轻GC的压力(./pool_test.go)
	- sync.Pool为每个P(gmp调度模型中的P)单独建立一个local缓存池,
		进一步降低高并发下对锁的竞争


	sync.Pool的一个典型应用就是建立像bytes.Buffer这样类型的临时
	缓存对象池:
	var bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	由于sync.Pool的Get方法从缓存池中挑选bytes.Buffer数据对象时并未考虑
	该数据对象是否满足调用者的需求, 因此一旦返回的Buffer对象是刚刚
	被"大数据"撑大后的, 并且即将被长期用于处理一些"小数据"时, 这个
	Buffer对象所占用的"大内存"将长时间得不到释放; 一旦这类情况集中
	出现, 将会给Go应用带来沉重的内存消耗负担;
	目前的go标准库采用两种方式来缓解这个问题:
	- 限制缓存池中的数据对象大小
	- 建立多级缓存池

*/

/*
	限制缓存池中的数据对象大小(TODO)

	$GOROOT/src/fmt/print.go
	// free saves used pp structs in ppFree; avoids an allocation per invocation.
	func (p *pp) free() {
		// Proper usage of a sync.Pool requires each entry to have approximately
		// the same memory cost. To obtain this property when the stored type
		// contains a variably-sized buffer, we add a hard limit on the maximum buffer
		// to place back in the pool.
		//
		// See https://golang.org/issue/23199
		if cap(p.buf) > 64<<10 {
			return
		}

		p.buf = p.buf[:0]
		p.arg = nil
		p.value = reflect.Value{}
		p.wrappedErr = nil
		ppFree.Put(p)
	}


	fmt包对于要放回缓存池的buffer对象做了一个限制性校验:
	如果buffer的容量大于64<<10, 则不让其回到缓存池中, 这样可以在一定程度
	上缓解处理小对象时重复利用大Buffer导致的内存占用问题
*/

/*
	建立多级缓冲池(TODO:源码)

	标准库的http包在处理http2数据时, 预先建立了多个不同大小的缓存池;
	// Buffer chunks are allocated from a pool to reduce pressure on GC.
	// The maximum wasted space per dataBuffer is 2x the largest size class,
	// which happens when the dataBuffer has multiple chunks and there is
	// one unread byte in both the first and last chunks. We use a few size
	// classes to minimize overheads for servers that typically receive very
	// small request bodies.
	//
	// TODO: Benchmark to determine if the pools are necessary. The GC may have
	// improved enough that we can instead allocate chunks like this:
	// make([]byte, max(16<<10, expectedBytesRemaining))
	var (
		http2dataChunkSizeClasses = []int{
			1 << 10,
			2 << 10,
			4 << 10,
			8 << 10,
			16 << 10,
		}
		http2dataChunkPools = [...]sync.Pool{
			{New: func() interface{} { return make([]byte, 1<<10) }},
			{New: func() interface{} { return make([]byte, 2<<10) }},
			{New: func() interface{} { return make([]byte, 4<<10) }},
			{New: func() interface{} { return make([]byte, 8<<10) }},
			{New: func() interface{} { return make([]byte, 16<<10) }},
		}
	)

	func http2getDataBufferChunk(size int64) []byte {
		i := 0
		for ; i < len(http2dataChunkSizeClasses)-1; i++ {
			if size <= int64(http2dataChunkSizeClasses[i]) {
				break
			}
		}
		return http2dataChunkPools[i].Get().([]byte)
	}

	func http2putDataBufferChunk(p []byte) {
		for i, n := range http2dataChunkSizeClasses {
			if len(p) == n {
				http2dataChunkPools[i].Put(p)
				return
			}
		}
		panic(fmt.Sprintf("unexpected buffer len=%v", len(p)))
	}

	以上可以根据要处理的数据的大小从最适合的缓存池中获取Buffer对象,
	并在完成数据处理后将对象归还到对应的池中, 而池中的所有临时buffer
	对象的容量始终是保持一致的, 从而尽量避免大材小用, 浪费内存的情况;

*/
