package main

/*
	将 ./non_table_drive_test.go 示例中重复的测试逻辑合并为一个, 并将
	预置的输入数据放入一个自定义结构体类型的切片中, 使之成为一个可扩展
	的测试设计, 即表驱动测试; 在编写Go测试代码时优先编写基于表驱动的测试;
*/

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	compareTests := []struct {
		a, b string
		i    int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", -1},
		{"x", "ab", 1},
		{"x", "a", 1},
		{"b", "x", -1},
	}

	/*
		无须改动后面的测试逻辑, 只需在切片中增加数据条目即可; 在这种测试
		设计中, 这个自定义结构体类型的切片(上述示例中的compareTests)就是
		一个表(自定义结构体类型的字段就是列), 而基于这个数据表的测试设计
		和实现则被称为"表驱动的测试";

		表驱动测试本身是编程语言无关的, 表驱动测试十分适合Go代码测试,
		go 团队在标准库和第三方项目中大量使用此种测试设计, 这样表驱动
		测试也就逐渐成为Go的一个惯用法;

		表驱动测试优点:
		- 简单和紧凑
		- 数据即测试: 表驱动测试的实质是数据驱动的测试, 扩展输入数据集
			即扩展测试; 通过扩展数据集, 可以很容易地实现提高被测目标
			测试覆盖率的目的
		- 结合子测试后, 可单独运行某个数据项的测试

	*/

	for _, tt := range compareTests {
		cmp := strings.Compare(tt.a, tt.b)
		if cmp != tt.i {
			t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i,
				tt.a, tt.b, cmp)
		}
	}
}

// 结合子测试后, 可单独运行某个数据项的测试
func TestCompareSubTest(t *testing.T) {
	compareTests := []struct {
		name, a, b string
		i          int
	}{
		{`compareTwoEmptyString`, "", "", 0},
		{`compareSecondParamIsEmpty`, "", "a", 1},
		{`compareFirstParamIsEmpty`, "a", "", -1},
		// 将测试结果的判定逻辑放入一个单独的子测试中,
		// 这样可以单独执行表中某项数据的测试
	}

	for _, tt := range compareTests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := strings.Compare(tt.a, tt.b)
			if cmp != tt.i {
				t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i,
					tt.a, tt.b, cmp)
			}
		})

	}
}

// go test -v -run /TwoEmptyString table_drive_test.go
// === RUN   TestCompare
// --- PASS: TestCompare (0.00s)
// === RUN   TestCompareSubTest
// === RUN   TestCompareSubTest/compareTwoEmptyString
// --- PASS: TestCompareSubTest (0.00s)
//     --- PASS: TestCompareSubTest/compareTwoEmptyString (0.00s)
// PASS
// ok  	command-line-arguments	0.001s
