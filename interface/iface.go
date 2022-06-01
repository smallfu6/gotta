package main

/*
	带方法签名的接口在运行时的具体结构由 iface 构成
	type iface struct {
		tab *itab
		data unsafe.Pointer
	}
	data: 存储了接口中动态类型的数据指针

	type itab struct {
		inter *interfacetype    // 接口的类型
		_type *_type			// 接口中存储的动态类型
		hash  uint32			// 接口的hash值
		_     [4]byte
		func [1]uintptr			// 动态数据类型的函数指针
	}

	inter 字段表示接口本身的类型, 类型 interfacetype 是对 _type 的简单包装:
	type interfacetype struct {
		typ		_type
		pkgpath name  // pkgPath?
		mhdr    []imethod
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
*/
