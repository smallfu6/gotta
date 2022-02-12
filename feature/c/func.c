

int my_function(int arg1, int arg2) {
	return arg1 + arg2;
}

int main() {
	int i = my_function(1, 2);
}

// gcc 和 clang 编译相同的 c 语言代码可能会生成不同的汇编指令, 不过生成的代码
// 在结构上不会有太大的区别; 这里使用 gcc 进行编译
// cc -S func.c  // 编译生成 func.s 汇编代码
