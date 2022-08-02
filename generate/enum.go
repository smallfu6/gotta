package main

/*
	go generate驱动生成枚举类型的String方法
*/

type Weekday int

// 利用自定义类型、const与iota可以模拟实现枚举常量类型
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// 通常会为Weekday类型手写String方法, 这样在打印上面枚举常量时能输出有意义的内容
func (d Weekday) String() string {
	/*
		TODO: stringer 工具
		如果一个项目中枚举常量类型有很多, 手动维护String方法十分烦琐且易错;
		对于这种情况, 使用go generate驱动stringer工具为这些枚举类型自动生成
		String方法的实现不失为一个较为理想的方案; ./stringer-demo

	*/
	switch d {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	}
	return "Sunday" // default 0 -> "Sunday"
}
