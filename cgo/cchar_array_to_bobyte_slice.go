package main

/*

	C语言中的数组与C中的指针在大部分场合可以随意切换; Go语言中的数组是原生
	的值类型, Go仅提供了C.GoBytes将C中的char类型数组转换为Go中的[]byte切片类型;

	而对于其他类型的C数组, 目前无法直接显式的在两者之间进行类型转换, 可以通过
	特定转换函数将C的特定类型数组转换为Go的切片类型(Go中数组是值类型, 其大小
	是静态的, 转换为切片更通用)
*/

// char cArray[] = {'a', 'b', 'c', 'd', 'e', 'f', 'g'};
// int intCArray[] = {1, 2, 3, 4, 5, 6, 7};
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	goArray := C.GoBytes(unsafe.Pointer(&C.cArray[0]), 7)
	// Go编译器并不能将C的cArray自动转换为数组的地址, 所以不能将数组变量直接
	// 传递给函数, 需要将数组第一个元素的地址传递给函数;
	fmt.Printf("%c\n", goArray) // [a b c d e f g]  %c 打印该值对应的unicode码值
	fmt.Printf("%v\n", goArray) // [97 98 99 100 101 102 103]

	goSlice := CArrayToGoArray(unsafe.Pointer(&C.intCArray[0]),
		unsafe.Sizeof(C.intCArray[0]), 7)
	fmt.Println(goSlice)

}

func CArrayToGoArray(cArray unsafe.Pointer, elemSize uintptr,
	length int) (goArray []int32) {
	for i := 0; i < length; i++ {
		j := *(*int32)((unsafe.Pointer(uintptr(cArray) + uintptr(i)*elemSize)))
		goArray = append(goArray, j)
	}
	return
}
