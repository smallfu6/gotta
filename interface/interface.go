package main

import "fmt"

/*
	go 语言中的接口是一组方法的签名, 一个接口类型的变量能够接收任何实现了此接口
	的用户自定义类型;  可以为任何自定义的类型添加方法, go 中没有任何形式的基于
	类型的继承, 而是使用接口来实现扁平化, 面向组合的设计模式;

	接口的优势:
		- 隐藏细节
		- 控制系统复杂性
		- 权限控制(接口限制了其调用者的行为, 降低安全风险)

	1. 隐式接口 ./implicit.go
	2. 指针和接口  ./pointer.go
	3. nil 和 no-nil

	接口动态类型:
	存储在接口中的类型称为接口中的类型称为接口的动态类型, 而接口本身的类型称
	为接口的静态类型;

	接口的比较性:(TODO:深入理解并熟练掌握)
		两个接口之间可以通过 == 或 != 进行比较;

		接口的比较规则:
		- 动态值为 nil 的接口变量总是相等(TODO: 不考虑动态变量, 验证)
		- 如果两个接口不为 nil 且接口变量有相同的动态类型和动态值, 那么这两个
			接口是相同的(TODO: 验证,要考虑动态类型是否可比较, 如果不能比较就
			是不同的接口?)
		- 如果接口存储的动态类型值是不可比较的, 则在运行时会报错(TODO:)

*/

type Shape interface {
	area() float64
}

type Rectangle struct {
	X, Y int
}

func (r Rectangle) area() float64 {
	return 0
}

/*
	类型断言:
	可以使用 i.(Type) 在运行时获取存储在接口中的类型, 编译时会保证类型 Type
	一定是实现了 i 的类型, 否则编译不通过; 同时还需要在运行时判断一次, 因为
	在类型断言方法 i.(Type) 中, 当 Type 类型实现了接口 i, 而接口内部没有任何
	动态类型(此时为nil)时, 在运行时就会 panic, 因为 nil 无法调用任何方法(TODO:?)

	为了避免运行时报错, 类型断言返回了第二个值 ok
	val, ok := i.(Type)
	可以通过判断返回的 bool 值判断接口变量i的动态类型是否是 Type;

	空接口:
	空接口增强了代码的扩展性和通用性, 获取空接口中动态类型的方法是:
	i.(type)
	type 是固定的关键字, 注意与带方法接口的断言的区别, 同时此语法只在 switch
	语句中有效
	switch f := arg.(type) {
	case bool:
		//...
	case string:
		//...
	}

*/
func main() {
	var s Shape
	rect := s.(Rectangle) // impossible type assertion: s.(Rectangle)
	fmt.Println(rect.area())
	// s.(type) // invalid AST: use of .(type) outside type switch
}
