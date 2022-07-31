package main

/*
	对于具名的C枚举类型xx, 可以通过C.enum_xx来访问该类型; 如果是匿名枚举,
	则只能访问其字段了;
*/

/*
 enum color {
		RED,
      BLUE,
      YELLOW,
 };
*/
import "C"
import "fmt"

func main() {
	var e, f, g C.enum_color = C.RED, C.BLUE, C.YELLOW
	fmt.Printf("%T\n", e) // uint32
	fmt.Println(e, f, g)  // 0, 1, 2
}
