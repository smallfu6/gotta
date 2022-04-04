package main

import "fmt"

/*
	- 结构体内嵌套匿名结构体时的初始化(√)
	- 结构体内嵌套自己不行, 可以嵌套自己的指针?
*/

type NestStruct struct {
	Attribute string
	Desktop   struct { // 结构体的定义不会被外部引用到
		IsWood   bool
		Location string
	}
}

func UseNestStruct() {
	var nest = NestStruct{
		Attribute: "woodDesk",
		// 在初始化这个被嵌入的结构体时，就需要再次声明结构才能赋予数据
		Desktop: struct {
			IsWood   bool
			Location string
		}{
			IsWood:   true,
			Location: "village",
		},
	}
	fmt.Println(nest)
}
