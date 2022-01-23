package main

/* 抽象语法树生成与类型检查 */

//  完成语法解析后，进入抽象语法树阶段, 在该阶段会将词法解析阶段生成的
//  AssignStmt 结构解析为一个 Node, Node 结构体是对抽象语法树中节点的抽象

// type Node struct {
// 	Left  *Node
// 	Right *Node
// 	Ninit Nodes
// 	Nbody Nodes
// 	List  Nodes
// 	Rlist Nodes
// 	Type  *types.Type
// 	E     interface{}
// 	...
// }
// 其中, Left(左节点)代表左边的变量 a, Right(右节点)代表整数333, 其 Op 操作为
// OLITERAL, Right 的 E 接口字段会存储值 333, 如果前一阶段为 IntLit 类型, 则
// 需要转换为 Mpint 类型, Mpint 类型用于存储整数常量, 结构如下:
// Mpint 代表整数常量
// type Mpint struct {
// 	Val  big.Int
// 	Ovf  bool
// 	Rune bool
// }
// 编译时 AST 阶段整数通过 math/big.Int 进行高精度存储, 浮点数通过 big.Float
// 进行高精度存储(TODO: math/big 源码)

// 在类型检查阶段, 右节点中的 Type 字段存储的类型会变为 types.Types[TINT],
// type.Types 是一个数组(var Types [NTYPE]*Type), 存储了不同标识对应的 go 语言
// 中的实际类型, 其中, types.Types[TINT] 对应 go 语言内置的 int 类型

// 接着完成最终的赋值操作, 并将右边常量的类型赋值给左边常量的类型
/*
TODO: $GOROOT/src/cmd/compile/internal/gc/typecheck.go, line3196

func typecheckas(n *Node) {
	if enableTrace && trace {
		defer tracePrint("typecheckas", n)(nil)
	}

	// delicate little dance.
	// the definition of n may refer to this assignment
	// as its definition, in which case it will call typecheckas.
	// in that case, do not call typecheck back, or it will cycle.
	// if the variable has a type (ntype) then typechecking
	// will not look at defn, so it is okay (and desirable,
	// so that the conversion below happens).
	n.Left = resolve(n.Left)

	if n.Left.Name == nil || n.Left.Name.Defn != n || n.Left.Name.Param.Ntype != nil {
		n.Left = typecheck(n.Left, ctxExpr|ctxAssign)
	}

	// Use ctxMultiOK so we can emit an "N variables but M values" error
	// to be consistent with typecheckas2 (#26616).
	n.Right = typecheck(n.Right, ctxExpr|ctxMultiOK)
	checkassign(n, n.Left)
	if n.Right != nil && n.Right.Type != nil {
		if n.Right.Type.IsFuncArgStruct() {
			yyerror("assignment mismatch: 1 variable but %v returns %d values", n.Right.Left, n.Right.Type.NumFields())
			// Multi-value RHS isn't actually valid for OAS; nil out
			// to indicate failed typechecking.
			n.Right.Type = nil
		} else if n.Left.Type != nil {
			n.Right = assignconv(n.Right, n.Left.Type, "assignment")
		}
	}

	if n.Left.Name != nil && n.Left.Name.Defn == n && n.Left.Name.Param.Ntype == nil {
		n.Right = defaultlit(n.Right, nil)
		n.Left.Type = n.Right.Type
	}

	// second half of dance.
	// now that right is done, typecheck the left
	// just to get it over with.  see dance above.
	n.SetTypecheck(1)

	if n.Left.Typecheck() == 0 {
		n.Left = typecheck(n.Left, ctxExpr|ctxAssign)
	}
	if !n.Left.isBlank() {
		checkwidth(n.Left.Type) // ensure width is calculated for backend
	}
}
*/

// TODO:
// 在 SSA 阶段再转换为 go 语言预置的标准类型(int, float64), 变量 a 中存储的
// 大数类型的 333 最终会调用 big.Int 包中的函数 Int64 函数并将其转换
// 为 int64 类型的常量, 形如: v4(?) == MOVQconst<int>[333](a[int])
