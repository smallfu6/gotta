package main

/*
	#include <stdio.h>
	union bar {
		char	c;
		int		i;
		double	d;
	};
*/
import "C"
import "fmt"

func main() {
	var b *C.union_bar = new(C.union_bar)
	fmt.Printf("%T\n", b) //  *[8]uint8
	fmt.Println(b)        //  &[0 0 0 0 0 0 0 0]
	// b.c = 4 // b.c undefined (type *[8]byte has no field or method c)
	// Go对待C的union类型与其他类型不同, Go将union类型看成[N]byte, 其中N为
	// union类型中最长字段的大小(圆整后的) TODO: double 占几个字节? N由最长
	// 字段决定, 如何存储下3个字段?
	b[0] = 4
	fmt.Println(b) //  &[4 0 0 0 0 0 0 0]
}
