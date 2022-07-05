package main

/*
	对没有初始化的channel(nil channel)进行读写操作将会发生阻塞; 但是
	nil channel 还有其他妙用
*/

import (
	"fmt"
	"time"
)

func main1() {
	c1, c2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(time.Second * 5)
		c1 <- 5
		close(c1)
	}()

	go func() {
		time.Sleep(time.Second * 7)
		c2 <- 7
		close(c2)
	}()

	var ok1, ok2 bool
	for {
		select {
		case x := <-c1:
			ok1 = true
			fmt.Println(x)
		case x := <-c2:
			ok2 = true
			fmt.Println(x)
		}

		if ok1 && ok2 {
			break
		}
	}
	fmt.Println("program end")
}

// 5
// 0
// 0
// 0
// 0
// 0
// 0
// ...
// 7
// program end
/*
	程序在输出5之后输出了很多0才输出7;
	第 5s, c1返回一个5后被关闭, select语句的case x := <-c1分支被选出执行,
	程序输出5, 回到for循环并开始新一轮select; c1被关闭, 由于从一个已关闭
	的channel接收数据将永远不会被阻塞, 所以新一轮select又将case x := <-c1
	这个分支选出并执行; c1处于关闭状态, 从这个channel获取数据会得到该
	channel对应类型的零值, 这里为0, 于是程序再次输出0; 程序按这个逻辑
	循环执行, 一直输出0值, 所以才会出现上述的输出;

	可以使用 nil channel 改进上述问题
*/

func main() {
	c1, c2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(time.Second * 5)
		c1 <- 5
		close(c1)
	}()

	go func() {
		time.Sleep(time.Second * 7)
		c2 <- 7
		close(c2)
	}()

	for {
		select {
		case x, ok := <-c1:
			if !ok {
				c1 = nil
			} else {
				fmt.Println(x)
			}
		case x, ok := <-c2:
			if !ok {
				c2 = nil
			} else {
				fmt.Println(x)
			}
			// 在判断c1或c2被关闭后, 显式地将c1或c2置为nil;对一个nil channel执行
			// 获取操作, 该操作将被阻塞, 因此已经被置为nil的c1或c2的分支将再也
			// 不会被select选中执行
		}
		if c1 == nil && c2 == nil {
			break
		}
	}
	fmt.Println("program end")
}
