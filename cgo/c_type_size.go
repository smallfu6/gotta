package main

/*
	Go提供了C.sizeof_T来获取C.T类型的大小, 如果是结构体、枚举及联合体类型,
	要在T前面分别加上struct_、enum_和union_的前缀;
*/

//	struct employee {
//		char *id;
//		int  age;
//	};
import "C"
import "fmt"

func main() {
	fmt.Printf("%#v\n", C.sizeof_int)
	fmt.Printf("%#v\n", C.sizeof_char)
	fmt.Printf("%#v\n", C.sizeof_struct_employee)
}
