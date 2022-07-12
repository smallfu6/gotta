package main

/*
	如果需要将所有测试函数放入一个更大范围的测试固件环境中执行, 就需要包
	级别测试固件; Go 1.4版本引入了TestMain, 使得包级别测试固件的创建和
	销毁得到支持;
*/

import (
	"fmt"
	"testing"
)

func setUp(testName string) func() {
	fmt.Printf("\tsetUp fixture for %s\n", testName)
	return func() {
		fmt.Printf("\ttearDown fixture for %s\n", testName)
	}
}

func TestFunc1(t *testing.T) {
	t.Cleanup(setUp(t.Name()))
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

func TestFunc2(t *testing.T) {
	t.Cleanup(setUp(t.Name()))
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

func TestFunc3(t *testing.T) {
	t.Cleanup(setUp(t.Name()))
	fmt.Printf("\tExecute test: %s\n", t.Name())
}

func pkgSetUp(pkgName string) func() {
	fmt.Printf("package SetUp fixture for %s\n", pkgName)
	return func() {
		fmt.Printf("package TearDown fixture for %s\n", pkgName)
	}
}

// TODO: testing.M
func TestMain(m *testing.M) {
	defer pkgSetUp("package demo_test")()
	m.Run()
}

// package SetUp fixture for package demo_test
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
// package TearDown fixture for package demo_test
// ok  	command-line-arguments	0.001s

// 可看到在所有测试函数运行之前, 包级别测试固件被创建; 在所有测试函数运行
// 完毕后, 包级别测试固件被销毁;
