package main

/*
	go 的 defer 会在当前函数返回前执行传入的函数, 经常用于关闭文件描述符,
	关闭数据库连接以及解锁资源;
	go defer 的实现一定是由编译器和运行时共同完成的;

	使用 defer 最常见的场景是在函数调用结束后完成一些收尾工作, 例如 defer
	回滚数据库的事务:
	func create(db *gorm.DB) error {
		tx := db.Begin()
		defer tx.Rollback()

		if err := tx.Create(&Post{Authot: "lujet"}).Error; err != nil {
			return err
		}

		return tx.Commit().Error
	}

	以上代码在创建了事务后, 就立刻调用 Rollback 保证事务一定会回滚; 即使事务
	执行成功了, 调用 tx.Commit 之后再执行 tx.Rollback 也不会影响已经提交的事务



	defer 的延迟调用机制:
	- 堆中分配: go1.1 ~ go1.12
		- 编译器将 defer 关键字转换为 runtime.deferproc, 并在调用 defer 关键字
			的函数返回之前插入 runtime.deferreturn
		- 运行时调用 runtime.deferproc 会将一个新的 runtime._defer 结构体追加
			到当前 goroutine 的链表头
		- 运行时调用 runtime.deferreturn 会从 goroutine 的链表中取出
			runtime._defer 结构并依次执行
	- 栈上分配: go1.13
		当该关键字在函数体中最多执行一次时, 编译期间的
			cmd/compile/internal/gc.state.call 会将结构体分配到栈上, 并调用
			runtime.deferprocStack
	- 开放编码: go.1.14 至今
		- 编译期间判断 defer 关键字, return 语句的数目确定是否开启开放编码优化
		- 通过deferBits 和  cmd/compile/internal/gc.openDeferInfo 存储 defer
			关键字信息
		- 如果 defer 关键字的执行可以在编译期间确定, 会在函数返回前直接插入相应
			代码, 否则会由运行时的 runtime.deferreturn 处理

	后调用 defer 的函数会先执行 defer:
		- 后调用的 defer 函数会被追加到 Goroutine_defer 链表的最前面(TODO)
		- 运行 runtime._defer 时是从前到后依次执行的
	会预先计算函数的参数:
		如果调用 runtime.deferproc 函数创建新的延迟调用, 就立刻复用函数的参数,
		函数的参数不会等到真正执行时计算;

*/

import (
	"fmt"
	"time"
)

//----------------------- defer 的调用时机以及多次调用时执行顺序的确定
// 作用域
func scope() {
	{
		defer fmt.Println("defer runs")
		fmt.Println("block ends")
	}
	fmt.Println("main ends")
}

// block ends
// main ends
// defer runs
/*
	defer 传入的函数不是在退出代码块的作用域时执行的, 它只会在当前函数的方法
	返回之前被调用(?)
*/

func scopeForManyDefer() {
	for i := 0; i < 5; i++ {
		defer fmt.Println(i) // 加入到当前 goroutine 的链表的最前面
	}
}

// 4
// 3
// 2
// 1
// 0
/*
	TODO: 后调用的先执行
*/

//------------ defer 使用传值的方式传递参数时会进行预计算, 会导致结果不符合预期
func preComputed() {
	startAt := time.Now()
	defer fmt.Println(time.Since(startAt))

	time.Sleep(time.Second)
}

// 103ns
// 向 defer 关键字传入匿名函数解决上面的问题:(其实是使用了匿名函数的延迟绑定)
// defer func() { fmt.Println(time.Since(startAt)) }()

/*
	TODO: 深入理解
	调用 defer 关键字会[立刻]复制函数中引用的外部参数(值传递), 所以
	time.Since(startAt) 的结果不是在 main 函数退出之前计算的,
	而是在 defer 关键字调用时计算的

	如果使用匿名函数, 利用闭包的延迟绑定, defer 后的函数真正执行时才会
	访问外部参数的值
*/
