//go:generate echo "htop"
package main

/*
	//go:generate command arg...
	注释符号//前面没有空格, 与go:generate之间亦无任何空格; 上面的go generate
	指示符可以放在Go源文件中的任意位置, 并且一个Go源文件中可以有多个go
	generate指示符, go generate命令会按其出现的顺序逐个识别和执行
*/

import "fmt"

//go:generate echo "middle"
func main() {
	fmt.Println("hello, go generate")

}

//go:generate echo "tail"

// htop
// middle
// tail
