package main

/*
	sync.Once 可以保证程序在运行器件某段代码只执行一次
	TODO: 应用场景
*/

import (
	"fmt"
	"sync"
)

func main() {
	o := &sync.Once{}
	for i := 0; i < 10; i++ {
		o.Do(func() {
			fmt.Println("only once")
		})
	}
}

/*
	 每一个 sync.Once 结构体中都包含一个用于标识代码块是否执行过的 done, 以及一个
	 互斥锁 sync.Mutex
	 type Once struct {
		done uint32
		m    Mutex
	 }
	 唯一暴露的方法 sync.Once.Do 接收一个入参为空的函数:
	 - 如果传入的函数已经执行过, 会直接返回;
	 - 如果传入的函数没有执行过, 会调用 sync.Once.doSlow 执行传入的参数;

	 func (o *Once) doSlow(f func()) {
		o.m.Lock()
		defer o.m.Unlock()
		if o.done == 0 {
			defer atomic.StoreUint32(&o.done, 1)
			f()
		}
	 }
	 - 为当前的 groutine 获取互斥锁;
	 - 执行传入的无入参函数;
	 - 运行延迟函数调用, 将成员变量 done 更新为1;
	 sync.Once 通过成员变量 done 确保函数不会执行第二次.

	 * 注意:
	 - sync.Once.Do 中传入的函数只会被执行一次, 哪怕函数中抛出了 panic
	 - 两次调用 sync.Once.Do 方法传入不同的函数只会执行第一次调用传入的函数;
*/
