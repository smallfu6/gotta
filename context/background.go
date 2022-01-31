package main

/*
	context 包中最常用的方法是 context.Background 和 context.TODO, 这两个方法都
	返回预先初始化好的私有变量 background  和 todo, 它们会在同一个程序中被复用.

	func Background()  Context {
		return background
	}

	func TODO() Context {
		return todo
	}

	var (
		background = new(emptyCtx)
		todo       = new(emptyCtx)
	)

	// An emptyCtx is never canceled, has no values, and has no deadline. It is not
	// struct{}, since vars of this type must have distinct addresses.
	type emptyCtx int

	func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
		return
	}

	func (*emptyCtx) Done() <-chan struct{} {
		return nil
	}

	func (*emptyCtx) Err() error {
		return nil
	}

	func (*emptyCtx) Value(key interface{}) interface{} {
		return nil
	}
	context.emptyCtx 通过空方法实现了 context.Context 接口中的所有方法, 没有
	任何功能.

	从代码来看, context.Background 和 context.TODO 互为别名, 没有太大差别,
	只是在语义上稍有不同:
	- context.Background 是上下文的默认值, 其他所有上下文都应该从它衍生出来
	- context.TODO 应该仅在不确定应该使用哪种上下文时使用

	多数情况下, 如果当前函数没有上下文作为入参, 会使用 context.Background
	作为起始上下文向下传递.


	Context 层级关系: (TODO: 理解)
                    ----------> Background ------> WithCancel
                   /
		Background -----------> WithCancel
		           \
				    \---------> WithValue
*/
