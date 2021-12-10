package main

/*
 * 在 ./closure.go 中闭包重写后, 需要遍历函数, 其核心逻辑在:
 * $GOROOT/src/cmd/compile/internal/gc/walk.go 文件的 walk 函数中.
 *
 * 该阶段会识别出声明但是并未被使用的变量, 遍历函数中的声明和表达式,
 * 将某些代表操作的节点转换为运行时的具体函数执行.
 */

// 获取 map 中的值会被转换为运行时 mapaccess2_fast64 函数(TODO)
// v, ok := m["foo"]
// 转换为
// autotmp_1, ok := runtime.mapaccess2_fast64(typeOf(m), m, "foo")
// v := *autotmp_1

// TODO
// 字符串变量的拼接会被转换为调用运行时 concatstrings 函数; 对于 new 操作,
// 如果发生了逃逸, 那么最终会调用运行时 newobject 函数将变量分配到堆区;
// for range 语句会重写为更简单的 for 语句形式.

// 在执行 walk 函数遍历之前, 编译器还需要对某些表达式和语句进行重新排序,
// 例如将 x /= y 替换为 x = x/y; 根据需要引入临时变量, 以确保形式简单,
// 例如 x = m[k] 或 m[k] = x, 而 k 可以寻址.
