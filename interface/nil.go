package main

import "fmt"

/* interface 类型不是任意类型 */

type TestStruct struct{}

func NilOrNot(v interface{}) bool {
	// TODO: 深入理解类型转换
	// 调用 NilOrNot 函数时会发生隐式类型转换, 除向方法传入参数外, 变量的赋值
	// 也会触发隐式类型转换; 在进行类型转换时, *TestStruct 类型会转换为
	// interface{} 类型, 转换后的变量不仅包含转换前的变量, 还包括变量的类型
	// 信息 TestStruct,  所以转换后的变量与 nil 不相等.
	return v == nil
}

func main() {
	var s *TestStruct
	fmt.Println(s == nil)    // true
	fmt.Println(NilOrNot(s)) // false
}
