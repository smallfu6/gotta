package main

import (
	"sync"
	"time"
)

/*
	goroutine和内核线程之间通过多路复用方式对应, 通常Go应用会启动很多goroutine,
	但创建的线程数量是有限的
*/

func goSleep() {
	time.Sleep(time.Second * 1000)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			goSleep()
			wg.Done()
		}()
	}
	wg.Wait()
}

// ps -ef | grep go_sleep
// lucas     765678   15917  0 16:50 pts/1    00:00:00 ./go_sleep
// cat /proc/765678/status | grep -i thread
// Threads:        6
// Speculation_Store_Bypass:       thread vulnerable

// 虽然额外启动了100个goroutine, 但进程使用的线程数仅为6, 这是因为Go优化
// 了一些原本会导致线程阻塞的系统调用, 比如time.Sleep及部分网络I/O操作,
// 通过运行时调度在不创建新线程的情况下依旧能达到同样的效果;
