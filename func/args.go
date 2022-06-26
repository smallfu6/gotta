package main

import "fmt"

/*
	函数的变长参数

	使用变长参数函数时最容易出现的一个问题是实参和形参不匹配
	对于参数是 "...T" 类型的形式参数, 其可以匹配和接受的实参类型有两种:
	- 多个 T 类型变量
	- t... (t 为 []T 类型变量)
	以上两种不可以同时混用作为函数的实参
*/

// dump 的形参: []interface{}
func dump(args ...interface{}) {
	for _, v := range args {
		fmt.Println(v)
	}
}

func main1() {
	s := []string{"Tony", "John", "Jim"}
	dump(s)
	// 可以认为只传入了一个实参 s, 其类型是切片实现了 interface,
	// 可以作为 []interface{} 的第一个元素

	// dump(s...)
	// cannot use s (variable of type []string) as []interface{} value
	// in argument to dump
	// s... 表明传入的切片类型, []string 类型变量并不能直接赋值给
	// 给 []interface{} 类型变量
}

func main() {
	// go 内置的 append 函数是一个例外, 其支持通过下面的方式将字符串附加到一个
	// 字节切片后面
	b := []byte{}
	b = append(b, "hello"...)
	fmt.Println(string(b))
	// string 类型本是不满足类型要求的(append 需要 []byte...), 但是 go 的编译器
	// 做了优化自动将string(只读的字节数组)隐式转换为了[]byte
}
