package main

import "fmt"

// 很多面向语言有接口概念, 例如 Java 和 C#, Java 的接口不仅可以定义方法签名,
// 还可以定义变量, 这些定义的变量可以直接在实现接口的类中使用

// Java 接口示例
// public interface MyInterface {
// 	public String hello = "hello";
// 	public void sayHello();
// }

// 实现 MyInterface 接口
// public class MyInterfaceImpl implements MyInterface {
// 	public void sayHello() {
// 		System.out.println (MyInterface.hello);
// 	}
// }

// Java 中的类必须通过上述方式显示声明实现的接口, 但是在 go 语言中, 实现接口
// 是隐式的; 定义接口需要使用 interface 关键字, 在接口中只能定义方法签名, 不能
// 包含成员变量, 如:
type error interface {
	Error() string
}

// 实现 error 接口
type RPCError struct {
	Code    int64
	Message string
}

func (e *RPCError) Eorror() string {
	return fmt.Sprintf("%s, code=%d", e.Message, e.Code)
}

// go 语言中接口的实现都是隐式的, 只需要实现 Error() string 方法, 就实现了
// error 接口; 在使用 RPCError 结构体时并不关心实现了哪些接口, go 语言只会在
// 传递参数, 返回参数以及变量赋值时检查类型是否实现了接口; 类型实现接口时只
// 需要实现接口中的全部方法.

/*
	类型:
		接口是 go 的内置类型, 它能够出现在变量的定义, 函数的入参和返回值中并
		对它们进行约束; 分为:
		runtime.iface 表示带有一组方法的接口,
		runtime.eface 表示不包含任何方法的接口 interface{}
		(TODO: $GOROOT/src/runtime/runtime2.go, line208)
*/

// interface{} 不是任意类型, 如果将某个类型转换为 interface{} 类型, 变量在运行
// 期间的类型也会发生变化, 获取变量类型时会得到 interface{}; TODO: 深入理解
