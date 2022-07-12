package main

import (
	"fmt"
	"testing"
)

/*
	一些测试函数所需的测试固件是相同的, 在 ./classic_package_level_testfixture_test.go
	的这种平铺模式下为每个测试函数都单独创建/销毁一次测试固件就显得有些重复
	和冗余; 在这样的情况下, 可以尝试采用测试套件来减少测试固件的重复创建;
*/

func suiteSetup(suiteName string) func() {
	fmt.Printf("\tsetUp fixture for suite %s\n", suiteName)
	return func() {
		fmt.Printf("\ttearDown fixture for suite %s\n", suiteName)
	}
}

func func1TestCase1(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func1TestCase2(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func1TestCase3(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func TestFunc1(t *testing.T) {
	t.Cleanup(suiteSetup(t.Name()))
	t.Run("testcase1", func1TestCase1)
	t.Run("testcase2", func1TestCase2)
	t.Run("testcase3", func1TestCase3)
}

func func2TestCase1(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func2TestCase2(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func func2TestCase3(t *testing.T) {
	fmt.Printf("\t\tExecute test: %s\n", t.Name())
}

func TestFunc2(t *testing.T) {
	t.Cleanup(suiteSetup(t.Name()))
	t.Run("testcase1", func2TestCase1)
	t.Run("testcase2", func2TestCase2)
	t.Run("testcase3", func2TestCase3)
}

func pkgSetUp(pkgName string) func() {
	fmt.Printf("package SetUp fixture for %s\n", pkgName)
	return func() {
		fmt.Printf("package TearDown fixture for %s\n", pkgName)
	}
}

func TestMain(m *testing.M) {
	defer pkgSetUp("package demo_test")()
	m.Run()
}

// package SetUp fixture for package demo_test
// === RUN   TestFunc1
// 	setUp fixture for suite TestFunc1
// === RUN   TestFunc1/testcase1
// 		Execute test: TestFunc1/testcase1
// === RUN   TestFunc1/testcase2
// 		Execute test: TestFunc1/testcase2
// === RUN   TestFunc1/testcase3
// 		Execute test: TestFunc1/testcase3
// 	tearDown fixture for suite TestFunc1
// --- PASS: TestFunc1 (0.00s)
//     --- PASS: TestFunc1/testcase1 (0.00s)
//     --- PASS: TestFunc1/testcase2 (0.00s)
//     --- PASS: TestFunc1/testcase3 (0.00s)
// === RUN   TestFunc2
// 	setUp fixture for suite TestFunc2
// === RUN   TestFunc2/testcase1
// 		Execute test: TestFunc2/testcase1
// === RUN   TestFunc2/testcase2
// 		Execute test: TestFunc2/testcase2
// === RUN   TestFunc2/testcase3
// 		Execute test: TestFunc2/testcase3
// 	tearDown fixture for suite TestFunc2
// --- PASS: TestFunc2 (0.00s)
//     --- PASS: TestFunc2/testcase1 (0.00s)
//     --- PASS: TestFunc2/testcase2 (0.00s)
//     --- PASS: TestFunc2/testcase3 (0.00s)
// PASS
