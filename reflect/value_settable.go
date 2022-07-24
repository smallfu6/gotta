package main

import (
	"fmt"
	"reflect"
)

/*
	reflect.Value 提供了 CanSet, CanAddr 及 CanInterface 等方法判断反射
	对象是否可设置(settable), 可寻址, 可恢复为一个 interface{} 类型变量
	TODO: 针对本节例子的实验结果, 深入研究引起不同类型的这些区别的底层
	原理

*/

type Person struct {
	Name string
	age  int
}

func main() {
	var n = 17
	// int
	fmt.Println("int:")
	val := reflect.ValueOf(n)
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // false false true

	// *int
	fmt.Println("\n*int:")
	val = reflect.ValueOf(&n)
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // false false true

	// *int Elem
	fmt.Println("\n*int to Elem:")
	val = reflect.ValueOf(&n).Elem()
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // true true true

	// *struct
	fmt.Println("\nptr to struct:")
	pp := &Person{"tony", 33}
	val = reflect.ValueOf(pp)
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // false false true

	// *struct Elem
	fmt.Println("\n*struct to Elem")
	val = val.Elem()
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // true true true

	// struct Field
	fmt.Println("\n*struct Elem Field to Name")
	val1 := val.Field(0) // Name
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val1.CanSet(), val1.CanAddr(), val1.CanInterface()) // true true true

	// struct Field for age
	fmt.Println("\n*struct Elem Field to age")
	val2 := val.Field(1) // age
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val2.CanSet(), val2.CanAddr(), val2.CanInterface()) // false true false
	// Person.age 未导出, 所以 CanInterface 为 false
	// 虽然 CanAddr 为 true, 但 CanSet 为 false

	// interface{}
	fmt.Println("\nstruct to interface")
	var i interface{} = &Person{"tony", 33}
	val = reflect.ValueOf(i)
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // false false true

	// interface{} Elem
	fmt.Println("\nstruct to interface Elem")
	val = val.Elem()
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // true true true

	// map

	// interface{}
	fmt.Println("\nmap")
	var m = map[string]int{
		"tony": 23,
		"jim":  34,
	}
	val = reflect.ValueOf(m)
	fmt.Printf("Settable = %v, CanAddr = %v, CanInterface = %v\n",
		val.CanSet(), val.CanAddr(), val.CanInterface()) // false false true

	val.SetMapIndex(reflect.ValueOf("tony"), reflect.ValueOf(12))
	fmt.Println(m) // map[jim:34 tony:12]
	/*
		map类型被反射对象比较特殊, key和value都是不可寻址和不可设置的;
		可以通过Value提供的SetMapIndex方法对map反射对象进行修改, 这种修改
		会同步到被反射的map变量中; TODO

		map元素不是一个变量, 不可以获取它的地址, 因为 map 的增长可能会导致
		已有的元素被重新散列到新的存储位置, 这样就可能会获取到无效的地址;
	*/

}

/*
	TODO: 深入研究 reflect 包后加深理解
	当被反射对象以值类型(T)传递给reflect.ValueOf时, 所得到的反射对象(Value)
	是不可设置和不可寻址的;
	当被反射对象以指针类型(*T或&T)传递给reflect.ValueOf时, 通过reflect.Value
	的Elem方法可以得到代表该指针所指内存对象的Value反射对象; 而这个反射对象
	是可设置和可寻址的, 对其进行修改(比如利用Value的SetInt方法)将会像函数
	的输出参数那样直接修改被反射对象所指向的内存空间的值;
	当传入结构体或数组指针时, 通过Field或Index方法得到的代表结构体字段
	或数组元素的Value反射对象也是可设置和可寻址的, 如果结构体中某个字段
	是非导出字段, 则该字段是可寻址但不可设置的(比如上面例子中的age字段);
	当被反射对象的静态类型是接口类型时(就像上面的interface{}类型变量i),
	该被反射对象的动态类型决定了其进入反射世界后的可设置性;
	如果动态类型为*T或&T时, 就像上面传给变量i的是&Person{}, 那么通过
	Elem方法获得的反射对象就是可设置和可寻址的;
*/
