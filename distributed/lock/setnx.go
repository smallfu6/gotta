package main

/*
	在分布式场景下, 可以使用 redis 的 setnx 命令实现高并发场景下对资源的抢占
	TODO: redis, setnx
*/

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func incr() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// TODO: 连接过程呢?

	var lockKey = "counter_lock"
	var counterKey = "counter"
	// lock
	resp := client.SetNX(lockKey, 1, time.Second*5)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		fmt.Println(err, "lock result: ", lockSuccess)
		return
	}

	// counter++
	getResp := client.Get(counterKey)
	cntValue, err := getResp.Int64()
	if err == nil {
		cntValue++
		resp := client.Set(counterKey, cntValue, 0)
		_, err := resp.Result()
		if err != nil {
			// log err
			println("set value error")
		}
	}

	println("curent counter is ", cntValue)
	delResp := client.Del(lockKey)
	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		println("unlock success!")
	} else {
		println("unlock failed", err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr()
		}()
	}
	wg.Wait()
}

/*
	通过代码和执行结果可以看到, 远程调用setnx实际上和单机的尝试锁非常相似,
	如果获取锁失败, 相关的任务逻辑就不应该继续向前执行;
	setnx很适合在高并发场景下, 用来争抢一些"唯一"的资源; 例如交易撮合系统中
	卖家发起订单, 而多个买家会对其进行并发争抢; 这种场景没有办法依赖具体的
	时间来判断先后, 因为不管是用户设备的时间, 还是分布式场景下的各台机器的时间,
	都是没有办法在合并后保证正确的时序的; 哪怕是同一个机房的集群, 不同的机
	器的系统时间可能也会有细微的差别;

*/
