package main

import (
	"fmt"
	"reflect"
)

/*

	TODO: 前期只做简单了解和使用, 底层原理熟悉后再做深入研究
	反射在大多数应用和服务中并不常见, 但是很多框架依赖 go 语言的反射机制简化
	代码; 因为 go 语言的语法元素很少, 设计简单, 所以其表达能力不够强, 但是 go
	的 reflect 包能弥补语法上的一些劣势;

	reflect 实现了运行时的反射能力(TODO: 理解), 能让程序操作不同类型的对象,
	reflect 包中有两对非常重要的函数和类型, 如下:
		- reflect.TypeOf() 获取类型信息, 对应 reflect.Type
		- reflect.ValueOf 获取数据的运行时表示(?), 对应 reflect.Value

	反射包中的所有方法基本都是围绕 reflect.Type 和 reflect.Value 两个类型
	设计的, 可以通过 TypeOf 和 ValueOf 将普通变量(interface)转换为包中提供
	的 reflect.Type 和 reflect.Value


	运行时反射是程序在运行期间检查其自身结构的一种方式, 反射带来的灵活性是一把
	双刃剑; 反射作为一种元编程方式, 可以减少重复代码; 但是过度使用反射会使程序
	逻辑变得难以理解并且运行缓慢, 使用 reflect 要遵循以下三大法则:
	1. interface{} 变量可以转换成反射对象;
		使用 reflect.TypeOf 和 reflect.ValueOf 能够获取 go 语言中变量对应的反射
		对象, 一旦获取了反射对象, 就能够得到跟当前类型相关的数据和操作, 并且可
		以使用这些运行时获取的结构执行方法;
	2. 从反射对象可以获取 interface{} 变量;
		从反射对象到接口值的过程是从接口值到反射对象的镜面过程, 两个过程需要
		经过两次转换:
		- 从接口值到反射对象
			- 从基本类型到接口类型的类型转换(隐式转换)
			- 从接口类型到反射对象的转换
		- 从反射对象到接口值
			-  反射对象转换成接口类型
			- 通过显示类型转换变成原始类型
	3. 要修改反射对象, 其值必须可设置;
		由于 go 语言的函数调用都是传值的, 所以得到的反射对象跟原变量没有任何
		关系, 那么直接修改反射对象就无法改变原变量, 程序为了防止错误会抛出
		panic, 需要:
		- 调用 reflect.ValueOf 获取变量指针
		- 调用 reflect.Value.Elem 获取指针指向的变量
		- 调用 reflect.Value.Set[Type] 更新变量的值


	方法调用:
	./add.go


	reflect 包提供了多种能力, 包括如何使用反射来动态修改变量, 判断类型是
	否实现了某些接口以及动态调用方法等功能;(TODO:实践)


	反射底层原理:
	// TypeOf returns the reflection Type that represents the dynamic type of i.
	// If i is a nil interface value, TypeOf returns nil.
	func TypeOf(i any) Type {
		eface := *(*emptyInterface)(unsafe.Pointer(&i))
		return toType(eface.typ)
	}

	// emptyInterface is the header for an interface{} value.
	type emptyInterface struct {
		typ  *rtype
		word unsafe.Pointer
	}

	reflect.TypeOf 将传入的接口变量转换为底层的实际空接口 emptyInterface,
	并获取空接口的类型值; reflect.Type 实质上是空接口结构体中的 typ 字段,
	是 *rtype 类型, go 中的任何具体类型的底层结构都包含这一类型;

	reflect.ValueOf 函数的核心是调用了 unpackEface 函数, reflect.Value 包含了
	接口中存储的值及类型, 除此之外还包含了特殊的 flag 标志; (TODO: 了解)
	// unpackEface converts the empty interface i to a Value.
	func unpackEface(i any) Value {
		e := (*emptyInterface)(unsafe.Pointer(&i))
		// NOTE: don't read e.word until we know whether it is really a pointer or not.
		t := e.typ
		if t == nil {
			return Value{}
		}
		f := flag(t.Kind())
		if ifaceIndir(t) {
			f |= flagIndir
		}
		return Value{t, e.word, f}
	}
*/

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFunc() {
	fmt.Println("call ReflectCallFunc")
}

// 遍历结构体字段
func LoopStructField() {
	var user = User{Id: 1, Name: "json", Age: 23}
	getType := reflect.TypeOf(user)
	getValue := reflect.ValueOf(user)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Println(field.Name, field.Type, value)
	}

	/*
		reflect.Type 的 Field 方法主要用于获取结构体的元数据, 其返回的结构
		体 StructField 如下:
		// A StructField describes a single field in a struct.
		type StructField struct {
			// Name is the field name.
			Name string

			// PkgPath is the package path that qualifies a lower case (unexported)
			// field name. It is empty for upper case (exported) field names.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			PkgPath string

			Type      Type      // field type
			Tag       StructTag // field tag string
			Offset    uintptr   // offset within struct, in bytes
			Index     []int     // index sequence for Type.FieldByIndex
			Anonymous bool      // is an embedded field
		}

		reflect.Value 的 Field 方法主要返回 reflect.Value
	*/
}
