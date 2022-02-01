package main

/*
	context.WithCancel 函数能从 context.Context 中衍生出新的子上下文, 并返回
	用于取消该上下文的函数; 一旦执行返回的取消函数, 当前上下文及其子上下文都会
	被取消, 所有的 goroutine 都会同步收到这一取消信号.


	context 子树的取消: (TODO: 源码)
                /------------> Goroutine(context, cancel) ------> Goroutine
		       /
	Goroutine ---------------> Goroutine
				\
				 \-----------> Goroutine


	TODO: 源码
	// Canceling this context releases resources associated with it, so code should
	// call cancel as soon as the operations running in this Context complete.
	func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
		if parent == nil {
			panic("cannot create context from nil parent")
		}
		c := newCancelCtx(parent)
		propagateCancel(parent, &c)
		return &c, func() { c.cancel(true, Canceled) }
	}
	- context.newCancelCtx 将传入的上下文封装成私有结构体 context.cancelCtx;
	- context.propagateCancel 会构建父子上下文之间的关联, 当父上下文被取消后,
	子上下文也会被取消

	// propagateCancel arranges for child to be canceled when parent is.
	func propagateCancel(parent Context, child canceler) {
		done := parent.Done()
		if done == nil {
			return // parent is never canceled
		}

		select {
		case <-done:
			// parent is already canceled
			child.cancel(false, parent.Err())
			return
		default:
		}

		if p, ok := parentCancelCtx(parent); ok {
			p.mu.Lock()
			if p.err != nil {
				// parent has already been canceled
				child.cancel(false, p.err)
			} else {
				if p.children == nil {
					p.children = make(map[canceler]struct{})
				}
				p.children[child] = struct{}{}
			}
			p.mu.Unlock()
		} else {
			atomic.AddInt32(&goroutines, +1)
			go func() {
				select {
				case <-parent.Done():
					child.cancel(false, parent.Err())
				case <-child.Done():
				}
			}()
		}
	}
	context.propagateCancel 的作用是在 parent 和 child 之间同步取消和结束的信
	号, 保证在 parent 被取消时, child 也会收到对应的信号, 不会出现状态不一致的
	情况


	// cancel closes c.done, cancels each of c's children, and, if
	// removeFromParent is true, removes c from its parent's children.
	func (c *cancelCtx) cancel(removeFromParent bool, err error) {
		if err == nil {
			panic("context: internal error: missing cancel error")
		}
		c.mu.Lock()
		if c.err != nil {
			c.mu.Unlock()
			return // already canceled
		}
		c.err = err
		if c.done == nil {
			c.done = closedchan
		} else {
			close(c.done)
		}
		for child := range c.children {
			// NOTE: acquiring the child's lock while holding parent's lock.
			child.cancel(false, err)
		}
		c.children = nil
		c.mu.Unlock()

		if removeFromParent {
			removeChild(c.Context, c)
		}
	}
	context.cancelCtx.cancel 的作用是关闭上下文中的 Channel 并向所有子上下文同
	步取消信号.


	TODO: 函数实现
	context.WithDeadline 在创建 context.timerCtx 的过程中判断了父上下文的截止
	日期与当前日期, 并通过 time.AfterFunc 创建定时器, 当时候超过截止日期后会
	调用 context.timerCtx.cancel 同步取消信号.
	context.timerCtx 内部不仅通过嵌入 context.cancelCtx 结构体继承了相关变量和
	方法, 还通过持有的定时器 timer 和截止事件 deadline 实现了定时取消功能.


*/
