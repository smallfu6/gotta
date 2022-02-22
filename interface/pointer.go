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

func mainF() {
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
	- 对于 Cat{}, 意味着 Quack 方法会接收一个全新的 Cat{}(传参传递了结构体内容)

	go 实战中对于此的解释:
	接口的方法集定义了一组关联到给定类型的值或者指针的方法, 定义方法时使用的
	接收者的类型决定了这个方法是关联到值还是关联到指针, 还是两个都关联;
	对于方法集的规则:
		如果使用指针接收者来实现一个接口, 那么只有指向那个类型的指针才能够
		实现对应的接口; 如果使用值接收者实现一个接口, 那么那个类型的值和和指针
		都能实现对应的接口;
		之所以有这样的规则, 是因为编译器并不是总能自动获得一个值的地址, 如下

*/

//---------------------------------------------编译器不是总能获取一个值的地址
type duration int

func (d *duration) pretty() string {
	return fmt.Sprintf("Duration: %d", d)
}

func main() {
	duration(42).pretty() // pretty is not in method of duration
	// ./pointer.go:65:14: cannot call pointer method on duration(42)
	// ./pointer.go:65:14: cannot take the address of duration(42)
}
