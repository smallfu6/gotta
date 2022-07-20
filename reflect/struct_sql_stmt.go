package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"
)

/*
	反射让静态类型语言go在运行时具备了某种基于类型信息的动态特性, 利用这种特性
	fmt.Println 在无法提前获知传入参数的真正类型的情况下可以对其进行正确的格式
	化输出; json.Marshal 也通过这种特性对传入的任意结构体类型进行解构并正确生
	成对应的json文本;


*/

// 构建 SQL 查询语句的例子
// 一种ORM(Object Relational Mapping, 对象关系映射)风格的实现
func ConstructQueryStmt(obj interface{}) (stmt string, err error) {
	// 仅支持 struct 或 struct 指针类型
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		err = errors.New("only struct is supported")
		return
	}

	buffer := bytes.NewBufferString("")
	buffer.WriteString("SELECT ")

	if typ.NumField() == 0 {
		err = fmt.Errorf("the type [%s] has no fields", typ.Name())
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if i != 0 {
			buffer.WriteString(", ")
		}

		column := field.Name
		if tag := field.Tag.Get("orm"); tag != "" {
			column = tag
		}
		buffer.WriteString(column)
	}

	stmt = fmt.Sprintf("%s FROM %s", buffer.String(), typ.Name())
	return
}

type Product struct {
	ID        uint32
	Name      string
	Price     uint32
	LeftCount uint32 `orm:"left_count"`
	Batch     string `orm:"batch_number"`
	Updated   time.Time
}

type Person struct {
	ID      string
	Name    string
	Age     uint32
	Gender  string
	Addr    string `orm:"address"`
	Updated time.Time
}

func main() {
	stmt, err := ConstructQueryStmt(&Product{})
	if err != nil {
		fmt.Println("construct query stmt for Product error:", err)
		return
	}
	fmt.Println(stmt)

	stmt, err = ConstructQueryStmt(&Person{})
	if err != nil {
		fmt.Println("construct query stmt for Person rror:", err)
		return
	}
	fmt.Println(stmt)
}

/*
	ConstructQueryStmt通过反射获得传入的参数obj的类型信息, 包括(导出)字段数量、
	字段名、字段标签值等, 并根据这些类型信息生成SQL查询语句文本; 如果结构体
	字段带有orm标签, 该函数会使用标签值替代默认列名(字段名); 如果将
	ConstructQueryStmt包装成一个包的导出API, 那么它可以被放入任何Go应用中,
	并在运行时为传入的任意结构体类型实例生成对应的查询语句;
*/
