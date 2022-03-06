package main

import "fmt"

/*
	TODO: 源码, $GOROOT/src/runtime/string.go
	字符串是一块连续的内存空间, 实际上是由字符组成的数组; c 语言的字符串使用
	字符数组 char[] 表示, 数组会占用一块连续的内存空间, 而内存空间存储的字节
	共同组成了字符串;
	go 的字符串是只读的字节数组, 意味着字符串会分配到只读的内存空间(TODO),
	所以字符串不能通过下标或者其他形式改变其中的数据;
	go 不支持直接修改 string 类型变量的内存空间(TODO), 可以通过在 string 和
	[]byte 类型之间反复转换实现修改:
	- 将这段内存复制到堆中或栈中;
	- 将变量的类型转换成 []byte 后并修改字节数据;
	- 将修改后的字节数组转换回 string;


	go 字符串这种不可变的特性保证不会引用到意外发生改变的值(TODO), 因为 go 语言
	的字符串可以作为哈希表的键, 所以如果哈希表的键是可变的, 不仅会增加哈希表
	实现的复杂度, 还可能会影响哈希表的比较;
*/

//  如果代码中存在字符串, 编译器会将其标记成只读数据 SRODATA
func main() {
	str := "hello"
	fmt.Println([]byte(str))
}

// go tool compile -S string.go
/*
 go.string."hello" SRODATA dupok size=5
        0x0000 68 65 6c 6c 6f
*/

// 数据结构
// 每一个字符串在运行时都会使用 reflect.StringHeader 表示, 其中包含指向字节数组
// 和数组的大小
type StringHeader struct {
	Data uintptr // uintptr 仅仅是存储指针的整形树, 无指针类型的语义
	Len  int
}

// 切片的结构体
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// 与切片的结构体相比, 字符串只少了一个表示容量的 Cap 字段, 而正是因为切片在
// go 语言的运行时表示与字符串高度相似, 所以常说字符串是只读的切片类型;
// 字符串作为只读的类型, 并不会直接向字符串追加元素改变其本身的内存空间, 所有
// 在字符串上的写入操作都是通过复制实现的;

//---------------------------------类型转换-----------------------------------
// 在使用 go 语言解析和序列化Json等数据格式时, 经常需要将数据在  string 和
// []byte 之间来回转换, 类型转换由一定开销, runtime.slicebytetostring 等函数
// 经常出现在火焰图(flame graph, TODO), 成为程序的性能热点;

// 从字节数组到字符串的转换需要使用 runtime.slicebytetostring 函数, 该函数
// 在函数体中农会先处理两种比较常见的情况, 长度为0或1的字节数组
// func slicebytetostring(buf *tmpBuf, ptr *byte, n int) (str string) {
// 	if n == 0 {
// 		// Turns out to be a relatively common case.
// 		// Consider that you want to parse out data between parens in "foo()bar",
// 		// you find the indices and convert the subslice to string.
// 		return ""
// 	}
// 	if raceenabled {
// 		racereadrangepc(unsafe.Pointer(ptr),
// 			uintptr(n),
// 			getcallerpc(),
// 			funcPC(slicebytetostring))
// 	}
// 	if msanenabled {
// 		msanread(unsafe.Pointer(ptr), uintptr(n))
// 	}
// 	if n == 1 {
// 		p := unsafe.Pointer(&staticuint64s[*ptr])
// 		if sys.BigEndian {
// 			p = add(p, 7)
// 		}
// 		stringStructOf(&str).str = p
// 		stringStructOf(&str).len = 1
// 		return
// 	}

// 	var p unsafe.Pointer
// 	if buf != nil && n <= len(buf) {
// 		p = unsafe.Pointer(buf)
// 	} else {
// 		p = mallocgc(uintptr(n), nil, false)
// 	}
// 	stringStructOf(&str).str = p
// 	stringStructOf(&str).len = n
// 	memmove(p, unsafe.Pointer(ptr), uintptr(n))
// 	return
// }

// 无论从字节数组转换为字符串, 还是从字符串转换为字节数组, 内存复制导致的性能
// 损耗会随着字符串和 []byte 长度的增长而增长; 作为只读的数据类型, 我们无法
// 改变字符串本身的结构, 但是在做拼接和类型转换等操作时一定要注意性能损耗, 遇到
// 需要极致性能的场景一定要尽量减少类型转换的次数;
