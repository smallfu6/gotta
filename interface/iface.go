package main

/*
	带方法签名的接口在运行时的具体结构由 iface 构成
	type iface struct {
		tab *itab
		data unsafe.Pointer
	}
	data: 存储了接口中动态类型的数据指针
	tab: 存储了接口的类型, 接口中的动态数据类型, 动态数据类型的函数指针

	type itab struct {
		inter *interfacetype    // 接口的类型
		_type *_type			// 接口中存储的动态类型
		hash  uint32			// 接口的hash值
		_     [4]byte
		fun [1]uintptr			// 动态数据类型的函数指针
	}

	inter 字段表示接口本身的类型, 类型 interfacetype 是对 _type 的简单包装:
	type interfacetype struct {
		typ		_type
		pkgpath name  // 代表接口所在的包名
		mhdr    []imethod // 表示接口中暴露的方法在最终可执行文件中的名字和类型的偏移量???
		// 通过偏移量在运行时能通过 resolveNameOff 和 resolveTypeOff 函数
		// 快速找到方法的名字和类型
		// func resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer
		// func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer
	}


	go 的各种数据类型都是在 _type 字段的基础上通过增加额外字段来管理的, 如以下
	的切片和结构体:
	type  slicetype struct {
		typ  _type
		elem *_type
	}

	type structtype struct {
		typ	 	_type
		pkgPath name
		fields  []structfield
	}

	其中 _type 包含了类型的大小, 哈希, 标志及偏移量等元数据
	type _type struct {
		size       uintptr
		ptrdata    uintptr
		hash       uint32
		tflag      tflag
		align      uint8
		fieldAlign uint8
		kind       uint8
		equal      func(unsafe.Pointer, unsafe.Pointer) bool
		gcdata     *byte
		str
		ptrToThis typeOff
	}

	itab.hash 是接口动态类型的唯一标识, 是 _type 类型中 hash 的副本; 在接口类型
	断言时可以使用该字段快速判断接口动态类型与具体类型 _type 是否一致, itab 中
	一个空的4字节用于对齐, fun 字段代表接口动态类型中的函数指针列表, 用于运行时
	接口调用动态函数; 值只定义了大小为1的数组[1]uintptr, 是其存储的是函数指
	针列表的第一个元素, 当运行时, 可以通过首地址+偏移找到任意的函数指针;

*/

/*
	接口内存逃逸分析
	iface 结构中的 data 字段存储了接口中具体值的指针, 这是因为存储的数据可能
	很大也可能很小, 很难预料; 同时也表明存储在接口中的值必须能够获取其地址,
	所以平时分配在栈中的值一旦赋值给接口后, 会发生内存逃逸, 在堆区为其开辟
	内存;

	示例程序 ./escape.go 中未发生逃逸, 与成书时的go版本有关; TODO: 寻找相似的
	例子验证以上结论
*/
