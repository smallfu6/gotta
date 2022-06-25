package main

import (
	"fmt"
	"io"
)

/*
	go 语言的方法声明比函数声明多了一个参数, 称之为 receiver 参数, 其是方法和
	类型之间的纽带;

	func (receiver T/*T) MethodName() {
	}
	方法声明中的T称为 receiver 的基类型, 上述方法被绑定到类型 T 上;
*/

// receiver 参数的基类型本身不能是指针类型或接口类型; TODO: 为什么?
type MyInt *int

func (m MyInt) String() string {
	// invalid receiver type MyInt (pointer or interface type)
	return fmt.Sprintf("%d", *(*int)(m))

}

type MyReader io.Reader

func (r MyReader) Read(p []byte) (int, error) {
	// invalid receiver MyReader (pointer or interface type)
	return r.Read(p)
}
