package main

/*
	管道是Unix/Linux上一种典型的并发程序设计模式, 也是Unix崇尚"组合"设计哲学
	的具体体现; Go中没有定义管道, 但是Go语言缔造者们显然借鉴了Unix的设计哲学,
	在Go中引入了channel这种并发原语, 而channel原语使构建管道并发模式变得容易;

	在Go中管道模式被实现成了由channel连接的一条"数据流水线", 在该流水线中,
	每个数据处理环节都由一组功能相同的goroutine完成; 在每个数据处理环节,
	goroutine都要从数据输入channel获取前一个环节生产的数据, 然后对这些数据
	进行处理, 并将处理后的结果数据通过数据输出channel发往下一个环节;
*/

// 根据起始和步长生成整形数字
func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

// 过滤调寄数
func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

// 求平方
func square(in int) (int, bool) {
	return in * in, true
}

func spawn(f func(int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

func main1() {
	in := newNumGenerator(2, 18)
	out := spawn(square, spawn(filterOdd, in))
	for v := range out {
		println(v)
	}
}

// 管道模式具有良好的可扩展性, 如果要在上面示例代码的基础上在最开始处新增
// 一个处理环节, 比如过滤掉所有大于1000的数(filterNumOver1000), 可以像下面
// 代码这样扩展管道流水线

func filterNumOver1000(n int) (int, bool) {
	if n > 1000 {
		return 0, false
	}
	return n, true
}

func main() {
	in := newNumGenerator(2, 18)
	out := spawn(filterNumOver1000, (spawn(square, spawn(filterOdd, in))))
	for v := range out {
		println(v)
	}
}

/*
	两种基于管道模式的扩展模式
	- 扇出模式(fan-out)
		在某个处理环节, 多个功能相同的goroutine从同一个channel读取数据并处理,
		直到该channel关闭; 使用扇出模式可以在一组goroutine中均衡分配工作量,
		从而更均衡地利用CPU;
	- 扇入模式(fan-in)
		在某个处理环节, 处理程序面对不止一个输入channel, 把所有输入channel
		的数据汇聚到一个统一的输入channel, 然后处理程序再从这个channel中
		读取数据并处理, 直到该channel因所有输入channel关闭而关闭;
*/
