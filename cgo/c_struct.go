package main

/*
	#include <stdlib.h>

	struct employee {
		char *id;
		int age;
	};
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	id := C.CString("1247")
	fmt.Printf("%T\n", id)
	defer C.free(unsafe.Pointer(id))

	var p = C.struct_employee{
		id:  id,
		age: 21,
	}
	fmt.Printf("%T\n", p)
	fmt.Printf("%#v\n", p)
}

// *main._Ctype_char
// main._Ctype_struct_employee
// main._Ctype_struct_employee{id:(*main._Ctype_char)(0x210d820), age:21, _:[4]uint8{0x0, 0x0, 0x0, 0x0}}
