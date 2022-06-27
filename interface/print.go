package main

/*
	每个接口类型变量在运行时的表示都是由两部分组成, 这两种接口类型可以分别
	标记为 eface(_type, data)  和 iface(tab, data); 虽然 eface 和 iface 的
	第一个字段不同, 但是 tab 和 _type 可统一看做动态类型的类型信息;

	go 语言中每种类型都有唯一的 _type 信息, 无论是内置原生类型, 还是自定义
	类型; go 运行时会为程序内的全部类型建立只读的共享 _type 信息表, 因此
	拥有相同动态类型的同类接口类型变量的 _type/tab 信息是相同的; 而接口类型
	变量的 data 部分则指向一个动态分配的内存空间, 该内存空间存储的是赋值
	给接口类型变量的动态类型变量的值;

	未显示初始化的接口类型变量的值为 nil,即该变量的 _type/tab 和 data 都为
	nil, 这样只需要判断两个接口类型变量的 _type/tab  是否相同和 data 指针
	所指向的内存空间存储的数据值是否相同(注意: 不是 data 指针的值)


	eface 和 iface 是 runtime 包中的非导出结构体定义, 不能直接在包外使用;
	提供了 println 预定义函数, 可以输出 eface 或 iface 的两个指针字段的值;
	println 在编译阶段会由编译器根据要输出的参数的类型将 println 替换为特
	定的函数, 这些函数定义在 $GOROOT/src/runtime/print.go 中:
	TODO: println 预定义函数

	func printeface(e eface) {
		print("(", e._type, ",", e.data, ")")
	}
	func printiface(i iface) {
		print("(", i.tab, ",", i.data, ")")
	}

	printeface 和 printiface 会输出各自的两个指针字段的值
*/

// nil 接口变量
func printNilInterface() {
	// nil 接口变量
	var i interface{} // 空接口类型
	var err error     // 非空接口类型
	println(i)        // (0x0,0x0)
	println(err)      // (0x0,0x0)
	println("i=nil: ", i == nil)
	println("err=nil: ", err == nil)
	println("i=err: ", i == err)

	/*
		无论是空接口类型变量还是非空接口类型变量, 若变量值为 nil, 则其内部
		表示均为 (0x0, 0x0)
	*/
}

// 空接口类型变量
func printEmptyInterface() {
	var eif1 interface{}
	var eif2 interface{}
	var n, m int = 17, 18

	eif1 = n
	eif2 = m

	println("eif1: ", eif1) // _type 相同, data 不同
	println("eif2: ", eif2)
	println("eif1=eif2: ", eif1 == eif2) // false
	// eif1:  (0x459d20,0xc000046768)
	// eif2:  (0x459d20,0xc000046760)
	// eif1=eif2:  false

	eif2 = 17
	println("eif1: ", eif1) // 相同
	println("eif2: ", eif2)
	println("eif1=eif2: ", eif1 == eif2) // true
	// eif1:  (0x459d20,0xc000046768)
	// eif2:  (0x459d20,0x4778b0)
	// eif1=eif2:  true
	/*
		go 在创建 eface 时一般会为 data 重新分配内存空间, 将动态类型变量的
		值复制到这块内存空间, 并将 data 指针指向这块内存空间; 因此在多数情况
		下看到的 data 值是不同的; 如上 0xc000046768 和 0x4778b0 指向的内存
		空间的值都是 17, 显然是直接指向了一块实现创建好的静态数据区(TODO)
		TODO: 同时要注意到在对 eif2 赋值后, eif2 的 data 指针长度的变化,
		似乎是指向了内存的只读段?

		同时 go 优化了对 data 的分配, 并不是每次都分配新的内存空间, 对比上下
		两个 eif2, 其 data 值都是 0x4778b0, 在对 eif2 赋值后(eif2 = int64(17))
		并没有为 data 重新分配内存空间;
	*/

	eif2 = int64(17)
	println("eif1: ", eif1) // _type 不同
	println("eif2: ", eif2)
	println("eif1=eif2: ", eif1 == eif2) // false
	println("")
	// eif1:  (0x459d20,0xc000046768)
	// eif2:  (0x459de0,0x4778b0)
	// eif1=eif2:  false

	/*
		对于空接口类型的变量, 只有在 eface._type 和 eface.data 所指向的数据的
		内容一致的情况下, 两个空接口类型变量才能相等;
		注意: eface.data 指针值不一定一致, 但是其指向的值必须一致
	*/
}

func main() {
	// printNilInterface()
	printEmptyInterface()
}
