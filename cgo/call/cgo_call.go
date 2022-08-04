package main

// void foo() {
// }
import "C"

// import "C"不支持放在xx_test.go文件中, 所以在这里进行调用, 在 test 中调用
// go 函数
func CallCFunc() {
	C.foo()
}

func foo() {
}

func CallGoFunc() {
	foo()
}
