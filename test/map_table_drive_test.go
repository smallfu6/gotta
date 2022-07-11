package main

/*
	表也可以使用基于自定义结构体的其他集合类型(如map)来实现
*/

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	// 用map来实现测试数据表
	compareTests := map[string]struct {
		a, b string
		i    int
	}{
		`compareTwoEmptyString`:     {"", "", 0},
		`compareSecondParamIsEmpty`: {"a", "", 1},
		`compareFirstParamIsEmpty`:  {"", "a", -1},
	}

	for name, tt := range compareTests {
		t.Run(name, func(t *testing.T) {
			cmp := strings.Compare(tt.a, tt.b)
			if cmp != tt.i {
				t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i,
					tt.a, tt.b, cmp)
				/*
					在表驱动的测试中, 数据表中的所有表项共享同一个测试结果的判定逻辑;
					这样需要在Errorf和Fatalf中选择一个来作为测试失败信息的输出途径,
					Errorf不会中断当前的goroutine的执行, 即便某个数据项导致了测试失败,
					测试依旧会继续执行下去, 而Fatalf恰好相反, 它会终止测试执行; 可根据
					测试的情况进行选择使用何种输出途径;

				*/
			}
		})
	}
}

// === RUN   TestCompare
// === RUN   TestCompare/compareTwoEmptyString
// === RUN   TestCompare/compareSecondParamIsEmpty
// === RUN   TestCompare/compareFirstParamIsEmpty
// --- PASS: TestCompare (0.00s)
//     --- PASS: TestCompare/compareTwoEmptyString (0.00s)
//     --- PASS: TestCompare/compareSecondParamIsEmpty (0.00s)
//     --- PASS: TestCompare/compareFirstParamIsEmpty (0.00s)
// PASS
// ok  	command-line-arguments	0.001s
