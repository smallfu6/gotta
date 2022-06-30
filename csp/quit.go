package main

import (
	"errors"
	"sync"
	"time"
)

/*
	退出模式
	程序中要启动多个 goroutine 协作完成应用的业务逻辑, 但是 goroutine 的运行
	形态很可能不同, 有的扮演服务端, 有的扮演客户端等等, 因此很难用一种统一
	的框架全面管理它们的启动、运行和退出; 可以把问题聚焦在实现一个"超时等待
	退出"框架, 以统一解决各种运行形态goroutine的优雅退出问题;


	一组goroutine的退出总体上有两种情况:
	- 并发退出: 在这类退出方式下, 各个goroutine的退出先后次序对数据处理
		无影响, 因此各个goroutine可以并发执行退出逻辑
	- 串行退出: 即各个goroutine之间的退出是按照一定次序逐个进行的, 次序若错了
		可能会导致程序的状态混乱和错误
*/

// 凡是实现了该接口的类型均可在程序退出时得到退出的通知和调用, 从而有机会做
// 退出前的最后清理工作
type GracefullyShutdowner interface {
	Shutdown(waitTimeout time.Duration) error
}

type ShutdownerFunc func(time.Duration) error

func (f ShutdownerFunc) Shutdown(waitTimeout time.Duration) error {
	return f(waitTimeout)
}

func ConcurrentShutdown(waitTimeout time.Duration,
	shutdowners ...GracefullyShutdowner) error {

	c := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown(waitTimeout)
			}(g)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()

	select {
	case <-c:
		return nil
	case <-timer.C:
		return errors.New("wait timeout")
	}
}

func SequentialShutdown(waitTimeout time.Duration,
	shutdowners ...GracefullyShutdowner) error {

	start := time.Now()
	var left time.Duration
	timer := time.NewTimer(waitTimeout)

	for _, g := range shutdowners {
		elapsed := time.Since(start)
		left = waitTimeout - elapsed
		c := make(chan struct{})
		go func(shutdowner GracefullyShutdowner) {
			shutdowner.Shutdown(left)
			c <- struct{}{}
		}(g)

		timer.Reset(left)
		select {
		case <-c:
			//  继续执行
		case <-timer.C:
			return errors.New("wait timeout")
		}
	}
	return nil
}

func shutdownMaker(processTm int) func(time.Duration) error {
	return func(d time.Duration) error {
		// time.Duration 参数有何意义, 只是为了满足 ShutdownerFunc 类型的签名?
		time.Sleep(time.Second * time.Duration(processTm))
		return nil
	}
}
