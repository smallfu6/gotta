package main

import "fmt"

/*
	无论是使用 for 的经典三段式循环还是 for-range 循环, 其底层都具有相同的
	汇编代码, 使用 for-range 的控制结构最终也会被 go 编译器转换为普通的 for
	循环;
*/

func main() {
	isForeverLoop()
}
func isForeverLoop() {
	arr := []int{1, 2, 3}
	for _, v := range arr {
		arr = append(arr, v)
	}
	fmt.Println(arr) // [1, 2, 3, 1, 2, 3]

	// for 循环会在堆上(栈?)申请新的内存地址存储 arr, _, v 等变量, 其作用域在
	// for 循环内; go 是值传递的, 所以在 for 语句中的 arr 只是函数内的 arr
	// 变量的一个副本; 对于所有 range 循环,  go 都会在编译期间将原切片或者
	// 数组赋值给一个新变量, 在赋值过程中就发生了复制, 并且通过 len 关键字
	// 预先获取了切片的长度, 所以在循环中追加新元素不会改变循环执行的次数;
}

// 清空切片
func clearSlice() {
	arr := []int{1, 2, 3}
	for i := range arr {
		arr[i] = 0
	}
	fmt.Println(arr)
}

// 依次遍历切片或哈希表然后置为0, 很耗费性能; 因为数组, 切片和哈希表占用的
// 内存空间都是连续的, 所以最快的方法是直接清空这块内存中的内容;
// TODO:汇编, 从汇编代码可看出, 编译器会直接使用 runtime.memclrNoHeapPointers
// 清空切片中的数据

/*
	map 使用 range 遍历与上述 slice 遍历同理, 在遍历内部增加 map 的长度不会
	改变 for 循环的次数;

	map(TODO: 底层结构, 溢出桶, 扩容)
	在遍历哈希表时, 编译器会使用 runtime.mapiterinit 和 runtime.mapiternext
	两个运行时函数重写 for-range 循环
	- runtime.mapiterinit:
		初始化 runtime.hiter 结构体中的字段, 并通过 runtime.fastrand 生成一个
		随机数后以便可以随机选择一个遍历桶的起始位置; go 语言团队在设计哈希
		表的遍历时, 不想让使用者依赖固定的遍历顺序, 所以引入了随机数保证遍历
		的随机性(TODO), map 的底层结构是俩表数组
	- runtime.mapiternext:
		- 当待遍历的桶为空时, 选择需要遍历的新桶
		- 当不存在待遍历的桶时, 返回(nil, nil) 键值对并中止遍历
*/

// 遍历字符串时与遍历数组, 切片和哈希表相似, 只是在遍历时会获取字符串中索引
// 对应字节并将其转换为 rune(int32), 在遍历字符串时拿到的值都是 rune 类型的
// 变量(rune由于表示的范围很大, 所以能处理一切字符); 使用下标访问字符串中的
// 元素时得到的就是字节(字符串是只读的字节数组切片)
