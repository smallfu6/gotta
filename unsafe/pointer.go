package main

/*
	unsafe.Pointer 虽然有可以打破类型安全屏障的能力, 但同时 unsafe.Pointer 的
	使用也需要遵循 unsafe 文档中定义的6条安全使用模式
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
