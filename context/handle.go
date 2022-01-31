package main

/*
	使用 context.Context 同步信号
	多个 goroutine 同时订阅 ctx.Done() Channel 中的消息, 一旦接收到取消信号
	立刻停止当前正在执行的工作(TODO: 为何可以同时订阅, 不阻塞, 广播机制?)
*/

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}

// process request with 500ms
// main context deadline exceeded
