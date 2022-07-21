package main

import (
	"fmt"
	"reflect"
)

/*
	通过Value提供的Index方法, 可以获取到切片及数组类型元素所对应的Value
	对象值(通过Value对象值我们可以得到其值信息); 通过Value的MapRange、
	MapIndex等方法, 可以获取到map中的key和value对象所对应的Value对象值,
	有了Value对象, 就可以像获取简单原生类型的值信息那样获得这些元素的值信息;
	对于结构体类型, Value提供了Field系列方法;


*/

type Person struct {
	Name string
	Age  int
}

func main() {
	// array
	var sl = [3]int{5, 6} // 数组
	val := reflect.ValueOf(sl)
	typ := reflect.TypeOf(sl)
	fmt.Printf("[%d %d %d]\n",
		val.Index(0).Int(),
		val.Index(1).Int(),
		val.Index(2).Int()) // [5, 6]
	fmt.Println(typ.Kind(), typ.String()) // array []int

	// map
	var m = map[string]int{
		"tony": 1,
		"jim":  2,
		"john": 3,
	}
	val = reflect.ValueOf(m)
	typ = reflect.TypeOf(m)
	iter := val.MapRange()
	fmt.Printf("{")
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		fmt.Printf("%s:%d,", k.String(), v.Int())
	}
	fmt.Printf("}\n")
	fmt.Println(typ.Kind(), typ.String()) // map map[string]int

	// struct
	var p = Person{"tony", 23}
	val = reflect.ValueOf(p)
	typ = reflect.TypeOf(p)
	fmt.Printf("{%s, %d}\n", val.Field(0).String(), val.Field(1).Int())
	fmt.Println(typ.Kind(), typ.Name(), typ.String()) // struct Person main.Person

	// channel
	var ch = make(chan int, 1) // channel
	val = reflect.ValueOf(ch)
	typ = reflect.TypeOf(ch)
	ch <- 17
	v, ok := val.TryRecv() // TODO: TryRecv 不阻塞
	// If the receive cannot finish without blocking, x is the zero Value and ok is false.
	// TODO: 如何判断通道是否关闭
	if ok {
		fmt.Println(v.Int()) // 17
	}
	fmt.Println(typ.Kind(), typ.String()) // chan chan int
}
