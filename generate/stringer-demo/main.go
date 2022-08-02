package main

/*
	TODO: go generate 只能基于包运行? 熟悉 stringer 工具
	go generate main.go
	stringer: error: 0 packages found
	main.go:23: running "stringer": exit status 1

	定义 go.mod 文件后成功生成 ./weekday_string.go
*/

import "fmt"

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// 为枚举类型生成String方法的实现
//go:generate  stringer -type=Weekday
func main() {
	var d Weekday
	fmt.Println(d)
	fmt.Println(Weekday(1))
}
