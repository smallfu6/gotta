package main

/*
	Go支持多返回值而C并不支持, 因此当将C函数用在多返回值的Go调用中时,
	C的errno将作为函数返回值列表中最后那个error返回值返回
	TODO: C的errno

*/

//#include <stdlib.h>
//#include <stdio.h>
//#include <errno.h>
// int foo(int i) {
//	  errno = 0;
//    if ( i > 5 ) {
//		errno = 8;
//		return i-5;
//    } else {
//		return i;
//    }
//}
import "C"
import "fmt"

func main() {
	i, err := C.foo(C.int(8))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i)
	}
}

// go run c_errno.go
// exec format error
// exec format error就是errno为8时的错误描述信息, 可以在C运行时库的errno.h中
// 找到errno=8与这段描述信息的联系:
// #define	ENOEXEC		8   /* Exec format error */
