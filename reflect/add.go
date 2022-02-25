package main

// TODO: 源码, 应用场景
// 作为一门静态语言, 通过 reflect 包利用反射在运行期间执行方法比较复杂; 但这
// 也是在静态语言中使用动态特性需要付出的成本

import (
	"fmt"
	"reflect"
)

func Add(a, b int) int { return a + b }

func main() {
	v := reflect.ValueOf(Add)
	if v.Kind() != reflect.Func {
		return
	}

	t := v.Type()
	argv := make([]reflect.Value, t.NumIn())
	for i := range argv {
		if t.In(i).Kind() != reflect.Int {
			return
		}
		argv[i] = reflect.ValueOf(i)
	}

	result := v.Call(argv)
	if len(result) != 1 || result[0].Kind() != reflect.Int {
		return
	}
	fmt.Println(result[0].Int())
}
