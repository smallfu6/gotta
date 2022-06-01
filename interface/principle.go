package main

/* 接口实现原理 */

/*
	编译时判断类型是否实现接口的方法:
	通常如果类型的方法与接口的方法完全是无序的状态, 并且类型有 m 个方法, 接口
	声明了 n 个方法, 总的时间复杂度在最坏的情况是 O(m×n), 即需要分别遍历类型
	与接口中的方法; go 语言在编译时对此做出的优化是先将类型与接口中的方法进行
	相同规则的排序, 再讲对应的方法进行比较;


			接口                           类型
			funcA -----------------------> funcA
			funcB -----------------------> funcB
			funcC -----------------------> funcC
			funcD -----------------------> funcD
			funcE -----------------------> funcE

	类型的方法可能少于或多于接口的方法, 虽然方法可能不会在相应的位置, 但是有序
	规则保证了当 funB 在接口方法列表中的序号为 i 时, 其在类型的方法列表中的
	序号大于或等于 i ;
	根据接口的有序规则, 遍历方法列表, 并在类型对应方法的列表序号i后查找是否存在
	对应的方法; 如果查找不到, 则说明类型未实现该接口, 编译时报错;  由于同一个
	类型或接口的排序在整个编译时只会进行一次, 因此排序的消耗可以忽略不计; 排序
	后最坏的时间复杂度仅为 O(m+n);

	以上功能的代码实现在 implements 函数中
	($GOROOT/src/cmd/compile/internal/typecheck/subr.go, line741)


*/

// implements reports whether t implements the interface iface. t can be
// an interface, a type parameter, or a concrete type. If implements returns
// false, it stores a method of iface that is not implemented in *m. If the
// method name matches but the type is wrong, it additionally stores the type
// of the method (on t) in *samename.
// func implements(t, iface *types.Type, m, samename **types.Field, ptr *int) bool {
// 	t0 := t
// 	if t == nil {
// 		return false
// 	}

// 	if t.IsInterface() || t.IsTypeParam() {
// 		if t.IsTypeParam() {
// 			// If t is a simple type parameter T, its type and underlying is the same.
// 			// If t is a type definition:'type P[T any] T', its type is P[T] and its
// 			// underlying is T. Therefore we use 't.Underlying() != t' to distinguish them.
// 			if t.Underlying() != t {
// 				CalcMethods(t)
// 			} else {
// 				// A typeparam satisfies an interface if its type bound
// 				// has all the methods of that interface.
// 				t = t.Bound()
// 			}
// 		}
// 		i := 0
// 		tms := t.AllMethods().Slice()
// 		for _, im := range iface.AllMethods().Slice() {
// 			for i < len(tms) && tms[i].Sym != im.Sym {
// 				i++
// 			}
// 			if i == len(tms) {
// 				*m = im
// 				*samename = nil
// 				*ptr = 0
// 				return false
// 			}
// 			tm := tms[i]
// 			if !types.Identical(tm.Type, im.Type) {
// 				*m = im
// 				*samename = tm
// 				*ptr = 0
// 				return false
// 			}
// 		}

// 		return true
// 	}

// 	t = types.ReceiverBaseType(t)
// 	var tms []*types.Field
// 	if t != nil {
// 		CalcMethods(t)
// 		tms = t.AllMethods().Slice()
// 	}
// 	i := 0
//  // 在比较之前, 会分别对接口与类型的方法进行排序, 排序使用了Sort 函数,
//  // 会根据元素数量选择不同的排序方法, 因为go根据函数名和包名可以唯一确定
//  // 命名空间中的函数, 所以排序后的结果是唯一的; (TODO: Sort 函数源码)
// 	for _, im := range iface.AllMethods().Slice() {
// 		if im.Broke() {
// 			continue
// 		}
// 		for i < len(tms) && tms[i].Sym != im.Sym {
// 			i++
// 		}
// 		if i == len(tms) {
// 			*m = im
// 			*samename, _ = ifacelookdot(im.Sym, t, true)
// 			*ptr = 0
// 			return false
// 		}
// 		tm := tms[i]
// 		if tm.Nointerface() || !types.Identical(tm.Type, im.Type) {
// 			*m = im
// 			*samename = tm
// 			*ptr = 0
// 			return false
// 		}
// 		followptr := tm.Embedded == 2

// 		// if pointer receiver in method,
// 		// the method does not exist for value types.
// 		rcvr := tm.Type.Recv().Type
// 		if rcvr.IsPtr() && !t0.IsPtr() && !followptr && !types.IsInterfaceMethod(tm.Type) {
// 			if false && base.Flag.LowerR != 0 {
// 				base.Errorf("interface pointer mismatch")
// 			}

// 			*m = im
// 			*samename = nil
// 			*ptr = 1
// 			return false
// 		}
// 	}

// 	return true
// }
