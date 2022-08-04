package main

/*
 */

//#include <unistd.h>
//void  cgoSleep() { sleep(1000); }
import "C"
import (
	"sync"
)

func cgoSleep() {
	C.cgoSleep()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			cgoSleep()
			wg.Done()
		}()
	}
	wg.Wait()
}

// ps -ef | grep cgo_sleep
// lucas     974629   15917  0 17:12 pts/1    00:00:00 ./cgo_sleep
// cat /proc/974629/status | grep -i thread
// Threads:        103
// Speculation_Store_Bypass:       thread vulnerable

/*
	Go调度器无法掌控C世界, 新创建的goroutine得到调度后, 会执行C空间的sleep函数
	进入睡眠状态, 执行这段代码的线程(M)也随之挂起, 之后Go运行时调度代码只能创
	建新的线程以供其他没有绑定M的P上的goroutine使用, 于是100个新线程被创建了
	出来;
	在日常开发中, 很容易在C空间中写出导致线程阻塞的C代码, 这会使得Go应用进程
	内线程数量暴涨的可能性大增, 与Go承诺的轻量级并发有背离;
*/
