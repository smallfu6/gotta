package main

/*
	Go测试代码的一般逻辑, 那就是针对给定的输入数据, 比较被测函数/方法返回
	的实际结果值与预期值, 如有差异则通过testing包提供的相关函数输出差异信息



*/

import (
	"strings"
	"testing"
)

// 使用了三组预置的测试数据对目标函数strings.Compare进行测试
func TestCompare(t *testing.T) {
	// 为被测函数/方法传入预置的测试数据, 然后判断被测函数/方法的返回结果
	// 是否与预期一致, 如果不一致则测试代码逻辑进入带有testing.Errorf的分支
	var a, b string
	var i int

	a, b = "", ""
	i = 0
	cmp := strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}

	a, b = "a", ""
	i = 1
	cmp = strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}

	a, b = "", "a"
	i = -1
	cmp = strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}

}
