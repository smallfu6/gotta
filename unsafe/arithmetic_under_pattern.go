package main

import (
	"fmt"
	"unsafe"
)

type Foo struct {
	s string
	b int
	c float64
	d [10]int
}

func main() {
	var a = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var foo = Foo{
		s: "foo",
		b: 17,
		c: 3.1415,
		d: a,
	}

	// 访问数组a的第四个元素
	var p = unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + 3*unsafe.Sizeof(a[0]))
	fmt.Println(*(*int)(p)) // 4

	// 访问Foo结构体的字段c
	p = unsafe.Pointer(uintptr(unsafe.Pointer(&foo)) + unsafe.Offsetof(foo.c))
	fmt.Println(*(*float64)(p)) // 3.1415

	// 对数组a的第一个元素进行逐字节步进式检查
	fmt.Println(unsafe.Sizeof((*int)(nil)))
	for i := uintptr(0); i < unsafe.Sizeof((*int)(nil)); i++ {
		p = unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + i)
		if i == 2 {
			xp := (*byte)(p)
			*xp = byte(2)
		}
		fmt.Printf("0x%x %x\n", (*byte)(p), *(*byte)(p))
	}
	fmt.Println(a)

	//...
}

// 4
// 3.1415
// 8
// 0xc00010a000 1   ->  00000001     |
// 0xc00010a001 0   ->  00000000     |
// 0xc00010a002 2   ->  00000010     |
// 0xc00010a003 0                    |
// 0xc00010a004 0                    |
// 0xc00010a005 0                   \|/ 高地址
// 0xc00010a006 0
// 0xc00010a007 0
// [131073 2 3 4 5 6 7 8 9 10]

// 此处使用了小端(小尾)模式: 数据的高字节保存在内存的高地址
// 131073 -> 10 00000000 00000001

/*
	模式3注意事项:
	- Offsetof 的使用不要越界
		如果在以上例子中, p	指向的地址超出原数组a的边界, 访问这块内存区域
		是有风险的, 尤其是尝试去修改它的时候; TODO: 实践如果修改会遇到的情况
	- unsafe.Pointer -> uintptr -> unsafe.Pointer 的转换要在一个表达式中
		func NewArray() *[10]int {
			a := [10]int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
			return &a
		}

		func main() {
			a := uintptr(unsafe.Pointer(NewArray()))
			// 存在风险: 此时间空隙, GC可能随时回收掉NewArray()返回的数组实例
			// TODO: GC
			p := unsafe.Pointer(a + unsafe.Sizeof(int(0)))
			fmt.Printf("%d\n", *(*int)(p))
		}


		uintptr 仅是一个整形值, 无指针语义, 无法起到对象引用的效果, 无法阻止
		GC回收内存对象; 若两次转换在一个表达式中, 则go编译器会保证两次转换
		期间NewArray函数返回的数组对象的有效性;
*/
