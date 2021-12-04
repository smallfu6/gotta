package main

/*
 * 词法解析:
 * 在词法解析阶段, 编译器会扫描输入的 go 源文件, 并将其符号(token)化, 如 "+"
 * 和 "-" 操作符会被转换为 _IncOp;
 * token 实质上是用 iota 声明的整数, 定义在:
 * $GOROOT/src/cmd/compile/internal/syntax/tokens.go 中; 符号化保留了 go 语言
 * 中定义的符号, 可以识别出错误的拼写. 同时, 字符串被转换为整数后, 在后续的
 * 阶段中能够被更加高效的处理(TODO).
 * 代码中声明的标识符, 关键字和分隔符等字符串都可以转换为对应的符号.
 */

// a := b + c(12) 对表达式符号化:
// a  -----> _Name
// := -----> _Define
// b  -----> _Name
// +  -----> _IncOp
// c  -----> _Name
// (  -----> _Lparen
// 12 -----> _Literal
// )  -----> _Rparen

import (
	"fmt"
	"go/scanner"
	"go/token"
)

// go 标准库 go/scanner, go/token 提供了接口用于扫描源代码
// 模拟对文本文件的扫描
func main() {
	src := []byte("cos(x) + 2i*sin(x) // Euler")

	// 初始化 scanner
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	// 扫描
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}
}

// 输出:
// 1:1     IDENT   "cos"
// 1:4     (       ""
// 1:5     IDENT   "x"
// 1:6     )       ""
// 1:8     +       ""
// 1:10    IMAG    "2i"
// 1:12    *       ""
// 1:13    IDENT   "sin"
// 1:16    (       ""
// 1:17    IDENT   "x"
// 1:18    )       ""
// 1:20    ;       "\n"
// 1:20    COMMENT "// Euler"
