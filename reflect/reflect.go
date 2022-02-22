package main

/*
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
	2. 从反射对象可以获取 interface{} 边量;
	3. 要修改反射对象, 其值必须可设置;

*/
