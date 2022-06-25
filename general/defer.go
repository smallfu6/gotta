package main

/*
	使用 defer 让函数更简洁, 更健壮

	在很多的函数中需要申请一些资源并在函数退出前释放或关闭这些资源, 如文件
	描述符, 数据库链接以及锁; 函数的实现需要确保这些资源在函数退出时被及时
	正确的释放, 无论函数的执行流是按预期顺利进行还是出现错误提前退出; 程序
	对错误进行处理时不能遗漏对资源的释放, 尤其是有多个资源需要释放的时候,
	大大增加了开发人员的心智负担; 此外, 当待释放的资源较多时, 代码逻辑性
	将变得十分复杂, 程序可读性, 健壮性也随之下降; 同时如果函数中抛出 panic,
	不使用 defer 也无法捕获并尝试从 panic 中恢复;

	使用 defer 便可以解决上述问题, go 的 defer 会在当前函数返回前执行传入
	的函数, 经常用于关闭文件描述符, 关闭数据库连接以及解锁资源;
	在 go 中, 只有在函数和方法内部才能使用 defer; defer 关键字后面只能接函数
	或方法, 这些函数被称为 deferred 函数; defer 将它们注册到其所在 goroutine
	用于存放 deferred 函数的栈数据结构中, 这些 deferred 函数将在执行 defer 的
	函数退出前按后进先出(LIFO)的顺序调度执行;

	无论是执行到函数尾部返回, 还是在某个错误处理分支显示调用 return 返回, 抑或
	出现 panic, 已经存储到 deferred 函数栈中的函数都会被调度执行, 因此,
	deferred 函数是一个在任何情况下都可以为函数进行收尾工作的函数;

	deferred 函数虽然可以拦截绝大部分的 panic, 但无法拦截并恢复一些运行时之外
	的致命问题, 如通过 c 代码制造的奔溃, deferred 便无能为力; (TODO: 实验)

	对于自定义的函数或方法, defer 可以无条件支持, 但是对于有返回值的自定义函数
	和方法, 返回值会在 deferred 函数被调度执行的时候自动丢弃;
	go 语言中除了有自定义的函数或方法, 还有内置函数, 其中只有 close, copy,
	delete, print, recover 等函数可以作为 deferred 函数, 而 append, cap,
	len, make, new 等内置函数是不可以直接作为 deferred 函数的(这些函数的返回
	值被丢弃, 使用 defer 无实际意义); 当然可以使用一个匿名函数包裹这些函数来
	间接满足需求, 但这么做似乎的实际意义需要自己把握;


	defer 的实现由编译器和运行时共同完成(TODO: 源码)

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
		- 运行 runtime._defer 时是从前到后依次执行的(入栈?)
	会预先计算函数的参数:
		如果调用 runtime.deferproc 函数创建新的延迟调用, 就立刻复用函数的参数,
		函数的参数不会等到真正执行时计算;



	在资源

*/

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
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
	调用 defer 关键字会[立刻]复制函数中引用的外部参数(值传递), 所以
	time.Since(startAt) 的结果不是在 main 函数退出之前计算的,
	而是在 defer 关键字调用时计算的;

	如果使用匿名函数, 利用闭包的延迟绑定, defer 后的函数真正执行时才会
	访问外部参数的值

	defer 关键字后面的表达式是在将 deferred 函数注册到 deferred 函数栈的时候
	进行求值的;
*/
func foo1() {
	sl := []int{1, 2, 3}
	fmt.Println(unsafe.Pointer(&sl))

	defer func(a []int) {
		fmt.Println(a)
	}(sl)

	sl = []int{3, 2, 1}
	fmt.Println(unsafe.Pointer(&sl))
}

func foo2() {
	sl := []int{1, 2, 3}
	fmt.Println(unsafe.Pointer(&sl))
	sliceD := (*reflect.SliceHeader)(unsafe.Pointer(&sl))
	fmt.Println(sliceD)

	defer func(p *[]int) {
		fmt.Println(unsafe.Pointer(p))
		sliceD1 := (*reflect.SliceHeader)(unsafe.Pointer(p))
		fmt.Println(sliceD1)
		fmt.Println(*p)
	}(&sl)

	sl = []int{3, 2, 1}
	fmt.Println(unsafe.Pointer(&sl))
	sliceD1 := (*reflect.SliceHeader)(unsafe.Pointer(&sl))
	fmt.Println(sliceD1)
}
