package main

import (
	"fmt"
	"time"
)

/*
	从原理看, 令牌通模型就是对全局计数的加减法操作过程, 但使用计数需要自己加
	读写锁; 也可以使用带缓冲的通道实现简单的加令牌/取令牌操作, 以下使用通道
	模拟简单的令牌桶模型
*/

var fillInterval = time.Millisecond * 10
var capacity = 100
var tokenBucket = make(chan struct{}, capacity)

func main() {
	go fillToken()
	time.Sleep(time.Hour)
}

func fillToken() {
	ticker := time.NewTicker(fillInterval)
	for {
		select {
		case <-ticker.C:
			select {
			case tokenBucket <- struct{}{}:
			default:
			}
			fmt.Println("current token count:", len(tokenBucket), time.Now())
		}
	}
}

func TakeAvailable(block bool) bool {
	var tokenResult bool
	if block {
		select {
		case <-tokenBucket:
			tokenResult = true
		}
	} else {
		select {
		case <-tokenBucket:
			tokenResult = true
		default:
			tokenResult = false
		}
	}
	return tokenResult
}

/*
	令牌桶每隔一段固定的时候向桶中放令牌, 如果记上一次放令牌的时间为t1,
	当时的令牌数k1, 放令牌的时间间隔为ti, 每次向令牌桶中放x个令牌, 令牌
	桶的容量为cap, 现在调用 TakeAvailable 来取n个令牌, 将这个时刻记为t2,
	那么在t2时刻, 令牌桶中理论上的令牌数为:
	cur = k1 + ((t2 - t1)/ ti) * x
	cur = cur > cap ? cap : cur

	用两个时间点的时间差再结合其他参数, 理论上在取令牌之前就完全可以知道
	桶里有多少令牌了, 本节向通道填充令牌的操作fillToken 理论上没有必要,
	只要在每次获取令牌的时候, 再对令牌同中的令牌数进行简单计算, 就可以得到
	正确的令牌数; 在得到正确的令牌数之后, 再进行实际的Take操作就可, 即只需
	要对令牌数进行简单的减法即可, 记得加锁以保证并发安全;
	限流器 github.com/juju/ratelimit 就是使用这种方法实现的; TODO: 源码

*/
