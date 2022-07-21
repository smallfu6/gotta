package main

import (
	"fmt"
	"reflect"
)

/*
	reflect.Value.Interface()是reflect.ValueOf()的逆过程, 通过Interface
	方法我们可以将reflect.Value对象恢复成一个interface{}类型的变量值;
	这个过程实质是将reflect.Value中的类型信息和值信息重新打包成一个
	interface{}的内部表示; TODO: 源码实现

*/

func main() {
	var i = 5
	val := reflect.ValueOf(i)
	r := val.Interface().(int)
	fmt.Println(r)

	r = 6
	fmt.Println(i, r)

	val = reflect.ValueOf(&i)
	q := val.Interface().(*int)
	fmt.Printf("%p, %p, %d\n", &i, q, *q)
	*q = 7
	fmt.Println(i)
}

/*
	通过reflect.Value.Interface()函数重建后得到的新变量(如例子中的r)与
	原变量(如例子中的i)是两个不同的变量, 它们的唯一联系就是值相同;
	如果反射的对象是一个指针(如例子中的&i), 那么通过reflect.Value.Interface()
	得到的新变量(如例子中的q)也是一个指针, 且它所指的内存地址与原指针变量相同;
	通过新指针变量对所指内存值的修改会反映到原变量上(变量i的值由5变为7);
*/
