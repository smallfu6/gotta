package main

import "fmt"

/*
	go 语言中的函数是唯一一种基于特定输入, 实现特定任务并可反馈任务执行结果的
	代码块, 本质上可以说 go 程序就是一组函数的集合;
	go 语言的函数具有如下特点:
		- 以 func 关键字开头
		- 支持多返回值
		- 支持具名返回值
		- 支持递归调用
		- 支持同类型的可变参数
		- 支持 defer, 实现函数优雅返回

	函数在 go 中属于"一等公民", ""一等公民"的定义:
		如果一门编程语言对某种语言元素的创建和使用没有限制, 可以像对待值(value),
		一样对待这种语法元素, 那么可以称这种语法元素是这门编程语言的"一等公民";

	拥有"一等公民"待遇的语法元素可以存储在变量中, 可以作为参数传递给函数, 可以
	在函数内部创建并可以作为返回值从函数返回; 在动态类型语言中, 语言运行时还
	支持对"一等公民"类型的检查; 除此外, 函数还可以被放入数组, 切片或 map 等
	结构中, 可以像其他类型变量一样被赋值给 interface{}, 甚至可以建立元素为函数
	的 channel;


	函数式编程(TODO)
	柯里化函数(currying):
		函数柯里化是把接受多个参数的函数变换成接受一个单一参数(原函数的第一个参数)
		的函数, 并返回接受余下的参数和返回结果的新函数的技术; 此技术以逻辑学家
		Haskell Curry 命令;
	函子(functor):(TODO)
		- 函子本身是一个容器类型, 这个容器可以是切片, map 或者 channel
		- 该容器类型需要实现一个方法, 该方法接收一个函数参数, 并在容器的每个
			元素上应用那个函数, 得到一个新函子, 原函子容器内部的元素值不受影响;

*/

//---------------------------------------------------------------- 柯里化函数
func times(x, y int) int {
	return x * y
}

// 柯里化函数[利用了函数的两点性质: 1.可以在函数中定义函数并返回, 2.闭包]
// 将原来接受两个参数的函数 times 转换为接受一个参数的函数 partialTimes 的过程
// 闭包: 在函数内部定义的匿名函数, 并且允许该函数访问定义它的外部函数的作用域,
// 本质上, 闭包是将函数内部和函数外部连接起来的桥梁
func partialTimes(x int) func(int) int {
	return func(y int) int {
		return times(x, y)
	}
}

func main1() {
	timesTwo := partialTimes(2)
	timesThree := partialTimes(3)

	timesFour := partialTimes(4)
	fmt.Println(timesTwo(5))   // 10
	fmt.Println(timesThree(5)) // 15
	fmt.Println(timesFour(5))  // 20
}

//---------------------------------------------------------------- 函子

type IntSliceFunctor interface {
	Fmap(fn func(int) int) IntSliceFunctor
}

// 函子
type intSliceFuncorImpl struct {
	ints []int
}

func (isf intSliceFuncorImpl) Fmap(fn func(int) int) IntSliceFunctor {
	newInts := make([]int, len(isf.ints))
	for i, elt := range isf.ints {
		retInt := fn(elt)
		newInts[i] = retInt
	}

	return intSliceFuncorImpl{ints: newInts}
}

func NewIntSliceFunctor(sl []int) IntSliceFunctor {
	return intSliceFuncorImpl{ints: sl}
}

func main() {
	// 原切片
	intSlice := []int{1, 2, 3, 4}
	fmt.Printf("init a functor from int slice: %#v\n", intSlice)

	f := NewIntSliceFunctor(intSlice)
	fmt.Printf("original functor: %+v\n", f)

	mapperFunc1 := func(i int) int {
		return i + 10
	}

	mapped1 := f.Fmap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)

	mapperFunc2 := func(i int) int {
		return i * 3
	}

	mapped2 := mapped1.Fmap(mapperFunc2)
	fmt.Printf("mapped functor2: %+v\n", mapped2)

	fmt.Printf("original functor: %+v\n", f) // 原函子没有改变
	fmt.Printf("composite functor: %+v\n",
		f.Fmap(mapperFunc1).Fmap(mapperFunc2)) // TODO: gorm 的 Where 的实现
}

/*
	函子非常适合用来对容器集合元素进行批量同构(?)处理, 而且代码也比每次都对容器
	中的元素进行循环处理要优雅, 简洁许多; 但要在 go 中发挥函子的最大效能, 还需要
	配合泛型使用, 否则就需要为每一种容器类型都实现一套对应的 Functor 机制;
*/
