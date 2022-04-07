package main

import (
	"context"
	"fmt"
	// "sync"
	"time"
)

// context 的传播关系:  父 context 的退出会导致所有子 context 的退出
// 子 context 的退出不会影响父亲 context
func main() {
	ctx := context.Background()
	before := time.Now()
	preCtx, _ := context.WithTimeout(ctx, 500*time.Millisecond)
	// var wg sync.WaitGroup
	// wg.Add(1)
	go func() {
		// defer wg.Done()
		childCt, _ := context.WithTimeout(preCtx, 300*time.Millisecond)
		select {
		case <-childCt.Done():
			after := time.Now()
			fmt.Println("child during:", after.Sub(before).Milliseconds())
		}
	}()

	select {
	case <-preCtx.Done():
		after := time.Now()
		fmt.Println("pre during:", after.Sub(before).Milliseconds())
	}
	// wg.Wait()
}

// child during: 100
// pre during: 100
