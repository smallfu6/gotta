package main

/*
	string 类型变量是不可变的(immutable), 常规方法将一个string类型变量转换为
	[]byte类型, go 会为[]byte类型变量分配一块新内存, 并将 string 类型变量的
	值复制到这块新内存中;
	基于 unsafe 包实现的 String2Bytes 函数并不需要额外的内存:
	转换后的[]byte变量与输入参数中的string类型变量共享底层存储(注意, 依旧无法
	通过对返回的切片的修改来改变原字符串); 而将[]byte变量转换为string类型则
	更简单, 因为[]byte内部表示是一个三元组(ptr, len, cap), string的内部表示为
	一个二元组(ptr, len), 通过 unsafe.Pointer 将[]byte的内部表示重新解释为
	string的内部表示, 这就是 Bytes2String 的原理;
*/

import (
	"reflect"
	"unsafe"
)

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
