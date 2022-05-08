package main

/*
	使用上下文传值
	context 包中的 context.WithValue 能从父上下文中创建子上下文, 传值的子上下
	文使用 context.valueCtx 类型:

	func WithValue(parent Context, key, val interface{}) Context {
		if parent == nil {
			panic("cannot create context from nil parent")
		}
		if key == nil {
			panic("nil key")
		}
		if !reflectlite.TypeOf(key).Comparable() {
			panic("key is not comparable")
		}
		return &valueCtx{parent, key, val}
	}

	context.valueCtx 结构体会将除 Value 外的 Err, Deadline 等方法代理到父上下
	文中, 它只会响应 context.valueCtx.Value 方法
	type valueCtx struct {
		Context
		key, val interface{}
	}

	func(c *valueCtx) Value(key interface{}) interface{} {
		if c.key == key {
			return c.val
		}
		return c.Context.Value(key)
	}

	如果 context.valueCtx 中存储的键值对与 context.valueCtx.Value 方法中传入的
	参数不匹配, 就会从父上下文中查找该键对应的值, 直到某个父上下文中返回 nil
	或者查找到对应的值.

*/

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// 多个 goroutine  同时读
	value := ctx.Value("site")
	fmt.Println(value)
}

func main() {
	// 写数据的时候必须由父context创建子context
	ctx := context.WithValue(context.Background(), "site", "segmentfault")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}

	wg.Wait()
}

/*
	可以发现写数据的时候必须创建新的子 context, 是在主 goroutine 中进行的,
	而读数据是在 worker 中进行的，若读不到则在父 context 中读数据; worker 中
	传入的 ctx 需要先创建, 可以认为读和写是有前后关系的, 且存在多个 goroutine
	并发读, 不存在goroutine 并发读写的情况， 所以 context 是并发安全的;
*/
