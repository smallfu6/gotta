package main

func myFunction(a, b int) (int, int) {
	return a + b, a - b
}

func main() {
	myFunction(66, 77)
}

// go tool compile -S -N -l main.go > main.s
// 使用 -N -l 参数表示编译器不对汇编代码进行优化
