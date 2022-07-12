package main

import (
	"fmt"
	"testing"
)

/*
	在运行每个测试函数TestXxx时, 都会先通过setUp函数建立测试固件;
	并在defer函数中注册测试固件的销毁函数, 以保证在每个TestXxx执
	行完毕时为之建立的测试固件会被销毁, 使得各个测试函数之间的测试
	执行互不干扰;

*/

/*
	在setUp中返回匿名函数来实现tearDown的好处是可以在setUp中利用闭包
	特性在两个函数间共享一些变量, 避免了包级变量的使用;
*/
func setUp(testName string) func() {
	fmt.Printf("\tsetUp fixture for %s\n", testName)
	return func() {
		fmt.Printf("\ttearDown fixture for %s\n", testName)
	}
}

// 在 Go 1.14版本以前, 测试固件的setUp与tearDown一般是这么实现的
func TestFunc1(t *testing.T) {
	// 注册测试固件的销毁函数
	defer setUp(t.Name())() // TODO: defer 后的函数是立即求值的? 熟练掌握这种巧妙的用法
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

func TestFunc2(t *testing.T) {
	defer setUp(t.Name())()
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

func TestFunc3(t *testing.T) {
	defer setUp(t.Name())()
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

// go test -v classic_testfixture_test.go
// === RUN   TestFunc1
// 	setUp fixture for TestFunc1
// 	Execute test: TestFunc1
// 	tearDown fixture for TestFunc1
// --- PASS: TestFunc1 (0.00s)
// === RUN   TestFunc2
// 	setUp fixture for TestFunc2
// 	Execute test: TestFunc2
// 	tearDown fixture for TestFunc2
// --- PASS: TestFunc2 (0.00s)
// === RUN   TestFunc3
// 	setUp fixture for TestFunc3
// 	Execute test: TestFunc3
// 	tearDown fixture for TestFunc3
// --- PASS: TestFunc3 (0.00s)
// PASS
// ok  	command-line-arguments	0.001s

/*
	Go 1.14版本testing包增加了testing.Cleanup方法, 为测试固件的销毁提供了
	包级原生的支持, 如下:

	func setUp() func() {
		...
		return func() {
		}
	}

	func TestXxx(t *testing.T) {
		t.Cleanup(setUp())
		...
	}

*/
