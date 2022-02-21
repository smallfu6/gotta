package main

import "fmt"

/*
	接口在定义一组方法时没有对实现的接收者做限制, 所以既可以使用类型作为方法的
	接收者实现接口, 也可以使用指向类型的指针作为方法的接收者实现接口; 因为结构
	体类型和指针类型是不同的, 所以在实现接口时这两种类型不能画等号;
	虽然两种类型不同, 但是 go 编译器会在结构体类型和指针类型同时实现一个方法
	时报错 "method redeclared".

*/

// 在使用结构体指针实现接口, 使用结构体初始化变量调用接口方法时无法通过编译
type Duck interface {
	Quack()
}

type Cat struct{}

func (c *Cat) Quack() { // 使用结构体指针实现接口
	fmt.Println("Quack")
}

// func (c Cat) Quack() { // 使用结构体实现接口
// 	fmt.Println("Quack")
// }

func main() {
	var c Duck = Cat{}
	// cannot use Cat literal (type Cat) as type Duck in assignment:
	//	Cat does not implement Duck (Quack method has pointer receiver)

	// var c Duck = &Cat{}
	c.Quack()
}

/*
	go 在传递参数时都是传值的;
	- 对于 &Cat{}, 意味着复制一个新的 &Cat{} 指针, 该指针与原来的指针指向相同
	且唯一的结构体, 所以编译器可以隐式对变量解引用(dereference)获取指针指向的
	结构体
	- 对于 Cat{}, 意味着 Quack 方法会接收一个全新的 Cat{}

*/
