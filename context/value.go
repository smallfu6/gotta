package main

/*
	使用上下文传值
	context 包中的 context.WithValue 能从父上下文中创建子上下文, 传值的子上下
	文使用 context.valueCtx 类型:

	func WithValue(parent Context, key, val interface{}) Context {
		if parent == nil {
			panic("cannot create context from nil parent")
		}
		if key == nil {
			panic("nil key")
		}
		if !reflectlite.TypeOf(key).Comparable() {
			panic("key is not comparable")
		}
		return &valueCtx{parent, key, val}
	}

	context.valueCtx 结构体会将除 Value 外的 Err, Deadline 等方法代理到父上下
	文中, 它只会响应 context.valueCtx.Value 方法
	type valueCtx struct {
		Context
		key, val interface{}
	}

	func(c *valueCtx) Value(key interface{}) interface{} {
		if c.key == key {
			return c.val
		}
		return c.Context.Value(key)
	}

	如果 context.valueCtx 中存储的键值对与 context.valueCtx.Value 方法中传入的
	参数不匹配, 就会从父上下文中查找该键对应的值, 直到某个父上下文中返回 nil
	或者查找到对应的值.


*/
