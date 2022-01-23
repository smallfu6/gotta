package main

// 词法解析与语法分析阶段
// $GOROOT/src/cmd/compile/internal/syntax/tokens.go line:99
type LitKind uint8

// TODO(gri) With the 'i' (imaginary) suffix now permitted on integer
//           and floating-point numbers, having a single ImagLit does
//           not represent the literal kind well anymore. Remove it?
const (
	IntLit   LitKind = iota // 代表整数
	FloatLit                // 代表浮点数
	ImagLit                 // 代表复数
	RuneLit
	StringLit
)

// TODO: 源码
// 在词法解析阶段, 会将赋值语句右边的常量解析为一个未定义的类型; 如 IntLit,
// go 源代码采用 utf-8 的编码方式, 在进行词法解析时, 当遇到需要赋值的常量操作
// 时会逐个读取后面常量的 utf-8 字符, 字符串的首字符为", 数字的首字符
// 为 '0'-'9', 具体实现位于 $GOROOT/src/cmd/compile/internal/syntax/scanner.go

// If the scanner mode includes the directives (but not the comments)
// flag, only comments containing a //line, /*line, or //go: directive
// are reported, in the same way as regular comments.
/*
func (s *scanner) next() {
	nlsemi := s.nlsemi
	s.nlsemi = false

	// ...
	switch s.ch {
	case -1:
		if nlsemi {
			s.lit = "EOF"
			s.tok = _Semi
			break
		}
		s.tok = _EOF

	case '\n':
		s.nextch()
		s.lit = "newline"
		s.tok = _Semi

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number(false)

	case '"':
		s.stdString()
*/
// 以赋值语句 a := 333 为例, 完成词法解析与语法分析时, 此赋值语句
// 将以 AssignStmt 结构表示
// type AssignStmt struct {
// 	Op       Operator
// 	Lhs, Rhs Expr
// 	simpleStmt
// }
