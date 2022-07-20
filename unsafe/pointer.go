package main

/*
	unsafe.Pointer 虽然有可以打破类型安全屏障的能力, 但同时 unsafe.Pointer 的
	使用也需要遵循 unsafe 文档中定义的6条安全使用模式

	go核心团队一直在完善工具链, 加强对代码中unsafe使用安全性的检查;
	通过go vet可以检查unsafe.Pointer和uintptr之间的转换是否符合下述
	6种安全模式;

	Go 1.14编译器在-race和-msan命令行选型开启的情况下, 会执行-d=checkptr检查,
	即对unsafe.Pointer进行下面两项合规性检查:
	- 当将*T1类型按模式1通过unsafe.Pointer转换为*T2时, T2的内存地址对齐系数
		不能高于T1的对齐系数
		./align_under_pattern1.go
	- 做完指针运算后，转换后的unsafe.Pointer仍应指向原先的内存对象
		./unsafe_compiler_checkptr.go

*/

/*
	模式1: *T1 -> unsafe.Pointer >- *T2
	本质就是内存块的重解释: 将原本解释为T1类型的内存重写解释为T2类型, 这是
	unsafe.Pointer 突破go类型安全屏障的基本使用模式

	// Float64bits returns the IEEE 754 binary representation of f,
	// with the sign bit of f and the result in the same bit position,
	// and Float64bits(Float64frombits(x)) == x.
	func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }

	// Float64frombits returns the floating-point number corresponding
	// to the IEEE 754 binary representation b, with the sign bit of b
	// and the result in the same bit position.
	// Float64frombits(Float64bits(x)) == x.
	func Float64frombits(b uint64) float64 { return *(*float64)(unsafe.Pointer(&b)) }

	以上的内存块重解释实现的类型转换不等价于go语法层面的显式类型转换

	func main() {
		var f float64 = 3.1415
		var d1 = uint64(f)
		var d2 = math.Float64bits(f)
		fmt.Printf("d1 = %d, d2 = %d\n", d1, d2) // d1 = 3, d2 = 4614256447914709615
		fmt.Printf("d1 = %d, d2 = %d\n", d1, d2) // d1 = 3, d2 = 4614
	}
	显式类型转换是语义层面的转换, float64 转换为 uint64 实质是取浮点数的整数
	部分; 而基于 unsafe.Pointer 使用模式1实现的 math.Float64bits 则只是机械的
	将这块内存视为 uint64 类型, 它并不在乎这块内存原先存储的是什么类型的数据;

	同时按照模式1使用 unsafe.Pointer 时也需要关注内存对齐问题
	./align_under_patter1.go
*/

/*
	模式2: unsafe.Pointer -> uintptr
	转换后的 uintptr 类型变量不再转换回 unsafe.Pointer
*/

/*
	模式3: 模拟指针运算
	操作任意内存地址上的数据都需要指针运算, go 常规语法不支持指针运算,
	但可以使用 unsafe.Pointer 的第三种安全使用模式模拟指针运算, 即在
	一个表达式中将 unsafe.Pointer 转换为 uintptr 类型, 使用 uintptr 类型
	的值进行算术运算后再转换回 unsafe.Pointer

	var b T
	var p  = unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + offset)
	*(*T)(p) = ...

	此模式常用于访问结构体内字段或数组中的元素, 也常用于实现对某内存对象的
	步进式检查; ./arithmetic_under_pattern.go
*/

/*
	模式4: 调用syscall.Syscall系列函数时指针类型到uintptr类型参数的转换
	同模式3 ./arithmetic_under_pattern.go  转换要在同一个表达式内
*/

/*
	模式5: 将reflect.Value.Pointer或reflect.Value.UnsafeAddr转换为指针
	Go标准库的reflect包的Value类型有两个返回uintptr类型值的方法:
	func (v Value) Pointer() uintptr
	func (v Value) UnsafeAddr() uintptr
	根据reflect文档的描述, 这两个方法是面向高级用户的;  使用uintptr类型
	作为返回值目的是促使这两个方法的调用者显式导入unsafe包并通过
	unsafe.Pointer对返回值进行转换[目前(Go 1.14版本)Go官方已经不建议继续
	使用这两个方法了(这两个方法处于Deprecated状态);
*/

/*
	模式6: reflect.SliceHeader和reflect.StringHeader必须通过模式1构建
	reflect包的SliceHeader和StringHeader两个结构体分别代表着切片类型和
	string类型的内存表示;
	可以通过模式1的内存块重解释来构造这两个结构体类型的实例;
	./slice_string_header.go

	如果通过常规语法定义一个reflect.SliceHeader类型实例并赋值, 后续反向
	转换成*[]T时存在SliceHeader.Data的值对应的地址上的对象已经被回收的风险;
	./slice_string_header_wrong.go

*/
