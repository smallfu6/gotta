package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/*
	Go语言内存管理是基于垃圾回收的, 垃圾回收例程会定期执行; 如果一块内存
	没有被任何对象引用就会被垃圾回收器回收; 而对象引用是通过指针实现的,
	unsafe.Pointer和其他常规类型指针一样, 可以作为对象引用; 如果一个对象
	仍然被某个unsafe.Pointer变量引用着, 那么该对象是不会被垃圾回收的,
	但是uintptr并不是指针, 它仅仅是一个整型值, 即便它存储的是某个对象的
	内存地址, 它也不会被算作对该对象的引用;
*/

type Foo struct {
	name string
}

func finalizer(p *Foo) {
	fmt.Printf("Foo: [%s]被垃圾回收\n", p.name)
}

func NewFoo(name string) *Foo {
	var f Foo = Foo{
		name: name,
	}
	// 在实例上设置了finalizer, 便于直观看到该实例是否在程序运行过程
	// 中被垃圾回收了
	runtime.SetFinalizer(&f, finalizer) // TODO
	return &f
}

func allocLarge() *[1000000]uint64 {
	a := [1000000]uint64{}
	return &a
}

func main() {
	var p1 = uintptr(unsafe.Pointer(NewFoo("FooRefByUintptr")))
	var p2 = unsafe.Pointer(NewFoo("FooRefByPointer"))

	for i := 0; i < 10; i++ {
		// 在每次循环中, 都会通过调用allocLargeObject做一些内存分配工作
		allocLarge()

		q1 := (*Foo)(unsafe.Pointer(p1))
		fmt.Printf("object ref by uintptr: %+v\n", *q1)

		q2 := (*Foo)(p2)
		fmt.Printf("object ref by pointer: %+v\n", *q2)

		runtime.GC() // 显式调用runtime.GC触发垃圾回收 TODO
		time.Sleep(1 * time.Second)
	}
}

// 为了避免编译器对程序进行内联优化, 在运行时传入了-gcflags="-l"命令行选项
// TODO: 内联优化的影响
// go run -gcflags="-l" uintptr.go
// object ref by uintptr: {name:FooRefByUintptr}
// object ref by pointer: {name:FooRefByPointer}
// Foo: [FooRefByUintptr]被垃圾回收
// object ref by uintptr: {name:FooRefByUintptr}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:FooRefByUintptr}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}
// object ref by uintptr: {name:}
// object ref by pointer: {name:FooRefByPointer}

// uintptr"引用"的Foo实例(FooRefByUintptr)在程序运行起来后很快就被回收了,
// 而unsafe.Pointer引用的Foo实例(FooRefByPointer)的生命周期却持续到程序终止;
// FooRefByUintptr实例被回收后, p1变量值中存储的地址值已经失效, 上面的输
// 出结果也证实了这一点: 这个地址处的内存后续被重新利用了;
