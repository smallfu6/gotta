package main

/*

	计算机程序可依据其瓶颈分为磁盘IO瓶颈型、CPU计算瓶颈型和网络带宽瓶颈型,
	分布式场景下有时候外部系统也会导致自身瓶颈; Web系统打交道最多的是网络,
	无论是接收、解析用户请求、访问存储, 还是把响应数据返回给用户, 都是要
	通过网络的; 在没有epoll/kqueue之类的系统提供的IO多路复用接口之前,
	多个核心的现代计算机最头痛的是C10k(TODO)问题, C10k问题会导致计算机没有办法
	充分利用CPU来处理更多的用户连接, 进而没有办法通过优化程序提升CPU利用率来
	处理更多的请求;
	自从Linux实现了epoll(TODO), FreeBSD实现了kqueue, 可以借助内核提供的API
	轻松解决C10k问题, 也就是说, 如今如果你的程序主要是和网络打交道, 那么瓶颈
	一定在用户程序而不在操作系统内核; 随着时代的发展, 编程语言对这些系统调用
	又进一步进行了封装, 如今做应用层开发几乎不会在程序中看到epoll之类的字眼,
	大多数时候我们只需要聚焦在业务逻辑上;
	Go的net库针对不同平台封装了不同的系统调用API, http库又是构建在net库之上的,
	所以在Go语言中我们可以借助标准库很轻松地写出高性能的http; ./server.go

	无论哪种类型的服务, 在资源使用到极限的时候都会导致请求堆积、超时、
	系统hang死, 最终伤害到终端用户; 对分布式的Web服务来说, 瓶颈还不一定
	总在系统内部, 也有可能在外部; 非计算密集型的系统往往会在关系型数据库
	环节失守, 而这时候Web模块本身还远远未达到瓶颈; 不管服务瓶颈在哪里,
	最终要做的事情都是一样的, 即流量限制;(TODO)


	流量限制的手段有很多, 常见的有漏桶和令牌桶两种
	1. 漏桶指有一个一直装满了水的桶, 每隔固定的一段时间即向外漏一地水; 如果
	接到了这滴水, 就可以继续服务请求, 如果没有接到, 就需要等待下一滴水;
	2. 令牌桶指匀速向桶中添加令牌, 服务请求时需要从桶中获取令牌, 令牌的数目
	可以按照需要消耗的资源进行相应的调整; 如果没有令牌, 可以选择等待或者放弃;

	这两种方法看起来很像, 其区别是:
	漏桶流出的速率固定, 而令牌桶只要在桶中还有令牌, 就可以继续拿; 也就是说,
	令牌桶是允许一定程度的并发的; 同时, 令牌桶在桶中没有令牌的情况下也会退化
	为漏桶;

	实际应用中令牌桶应用较为广泛, 开源界流行的限流器大多数都是基于令牌桶思想
	的, 并且在此基础上进行了一定程度的扩充;
	如 github.com/juju/ratelimit 就提供了几种不同特色的令牌桶填充方式:

	func NewBucket(fillInterval time.Duration, capacity int64) *Bucket
	默认的令牌桶, fillInterval 指每过多长时间向桶里放一个令牌, capacity 是桶的
	容量, 超过桶容量的部分被直接丢弃, 桶初始是满的;

	func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket
	每次向桶中放令牌时, 是放quantum个令牌, 而不是一个令牌;

	func NewBucketWithRate(rate float64, capacity int64) *Bucket
	会按照提供的比例, 每秒钟填充令牌数; 例如 capacity 是100, rate是0.1,
	那么每秒会填充10个令牌;

	从桶中获取令牌也提供了以下API:
	func (tb *Bucket) Take(count int64) time.Duration {}
	func (tb *Bucket) TakeAvailable(count int64) int64 {}
	func (tb *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {}
	func (tb *Bucket) Wait(count int64) {}
	func (tb *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {}

	缺点是这个库不支持不支持令牌桶预热(TODO), 且无法修改初始的令牌容量, 所以
	个别极端情况下的需求无法满足;

	./bucket.go

*/
