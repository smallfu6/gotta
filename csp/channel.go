package main

import (
	"fmt"
	"time"
)

/*
	由于无缓冲channel的运行时层实现不带有缓冲区, 因此对无缓冲channel的接收
	和发送操作是同步的, 即对于同一个无缓冲channel, 只有在对其进行接收操作
	的goroutine和对其进行发送操作的goroutine都存在的情况下通信才能进行, 否
	则单方面的操作会让对应的goroutine陷入阻塞状态; 如果一个无缓冲channel
	没有任何goroutine对其进行接收操作, 一旦有goroutine先对其进行发送操作,
	那么动作发生和完成的时序如下:
		发送动作发生
		-> 接收动作发生(有goroutine对其进行接收操作)
		-> 发送动作完成/接收动作完成(先后顺序不能确定) TODO:
	如果一个无缓冲channel没有任何goroutine对其进行发送操作, 一旦有goroutine
	先对其进行接收操作, 那么动作发生和完成的时序如下:
		接收动作发生
		-> 发送动作发生(有goroutine对其进行发送操作)
		-> 发送动作完成/接收动作完成(先后顺序不确定) TODO:

	根据上述时序结果, 对于无缓冲channel而言, 得到以下结论:
		发送动作一定发生在接收动作完成之前
		接收动作一定发生在发送动作完成之前
		(TODO: 如何理解)


	len(channel) 的应用:
	len是Go语言原生内置的函数, 可以接受数组、切片、map、字符串或
	channel类型的参数, 并返回对应类型的"长度"——一个整型值;
	以len(s)为例:
	如果s是字符串(string)类型, len(s)返回字符串中的字节数;
	如果s是[n]T或*[n]T的数组类型, len(s)返回数组的长度n;
	如果s是[]T的切片(slice)类型, len(s)返回切片的当前长度;
	如果s是map[K]T的map类型, len(s)返回map中已定义的key的个数;
	如果s是chan T类型, 那么len(s)针对channel的类型不同:
	- 当s为无缓冲channel时, len(s)总是返回0;
	- 当s为带缓冲channel时, len(s)返回当前channel s中尚未被读取的元素个数;

*/

/*
	针对带缓冲channel的len调用, 是否可以使用len函数来实现带缓冲channel的
	"判满", "判有"和"判空"逻辑?

	var c chan T = make(chan T, capacity)

	// 判空
	if len(c) == 0 {
		// 此时channel c 为空?
	}

	// 判有
	if len(c) > 0 {
		// 此时channel c 有数据?
	}

	if len(channel) == cap(channel) {
		// 此时 channel c 满了?
	}

	channel原语用于多个goroutine间的通信, 一旦多个goroutine共同对channel进行
	收发操作, 那么len(channel)就会在多个goroutine间形成竞态, 单纯依靠
	len(channel)来判断channel中元素的状态, 不能保证在后续对channel进行收
	发时channel的状态不变;
					|									|
					|									|
					| goroutine1						| goroutine2
					|									|
					|									|
					|									|
				   \|/									|
		if len(channel) == 0							|
					|                                   |
					| No                                |
					|                                   |
	----------------|-----------------------------------|
					|								   \|/
					|        竞                  -------------------
	len(channel) =1 |        态                  |  从 channel 中   |
					|        窗                  |  读取数据        |
					|        口                  |------------------|
					|                                   |
	----------------|-----------------------------------|
					|                                   |
				   \|/                                  |
			-------------------------                   |
			| 从 channel 中读取数据 |					|
			|-----------------------|					|
                                                       \|/
			len(channel) = 0

	goroutine1在使用len(channel)判空后, 便尝试从channel中接收数据; 但在其真
	正从channel中读数据前, goroutine2已经将数据读了出去, goroutine1后面的读
	取将阻塞在channel上, 导致后面逻辑失效; 因此为了不阻塞在channel上, 常见
	的方法是将判空与读取放在一个事务中，将判满与写入放在一个事务中,
	而这类事务我们可以通过select实现;

*/

func producer(c chan<- int) {
	var i int = 1
	for {
		time.Sleep(2 * time.Second)
		ok := trySend(c, i)
		if ok {
			fmt.Printf("[producer]: send [%d] to channel\n", i)
			i++
			continue
		}
		fmt.Printf("[producer]: try send [%d], but  channel is full\n", i)
	}
}

// select 事务
func tryRecv(c <-chan int) (int, bool) {
	select {
	case i := <-c:
		return i, true
	default:
		return 0, false
	}
}

// select 事务
func trySend(c chan<- int, i int) bool {
	select {
	case c <- i:
		return true
	default:
		return false
	}
}

func consumer(c <-chan int) {
	for {
		i, ok := tryRecv(c)
		if !ok {
			fmt.Println("[consumer]: try to recv from channel, but the channel is empty")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("[consumer]: recv [%d] from channel\n", i)
		if i >= 3 {
			fmt.Println("[consumer]: exit")
			return
		}
	}
}

func main() {
	c := make(chan int, 3)
	go producer(c)
	go consumer(c)
	select {}
}

/*
	上述使用事务的这种方法适合大多数场合, 但是它改变了channel的
	状态: 接收或发送了一个元素; 有时想在不改变channel状态的前提下
	单纯地侦测channel的状态, 又不会因channel满或空阻塞在channel上;
	很遗憾目前没有一种方法既可以实现这样的功能又适用于所有场合,
	在特定的场景下可以用len(channel)来实现:
	- 多发送单接收的场景: 即有多个发送者, 但有且只有一个接收者; 可以在
		接收者goroutine中根据len(channel)是否大于0来判断channel中是否
		有数据需要接收;
	- 多接收单发送的场景m 即有多个接收者, 但有且只有一个发送者; 可以在
		发送goroutine中根据len(channel)是否小于cap(channel)来判断是否
		可以执行向channel的发送操作;
*/
