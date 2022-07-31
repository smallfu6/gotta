package main

// typedef int myint;
// struct employee {
// 	char *id;
// 	int age;
// };
// typedef struct employee myemployee;
import "C"
import "fmt"

func main() {
	var a C.myint = 5
	var b C.struct_myemployee
	// _Ctype_struct_myemployee is incomplete (or unallocatable); stack allocation disallowed
	fmt.Println(a, b)

	// var p = C.struct_myemployee{
	// 	id:  id,
	// 	age: 21,
	// }
	// TODO: 不能解析出 myemployee 的字段
}

/*
	对原生类型的别名, 直接访问这个新类型名即可; 而对于复合类型的别名, 需要
	根据原复合类型的访问方式对新别名进行访问, 比如myemployee的实际类型为struct,
	使用myemployee时也要加上struct_前缀;

*/
