

// 当 my_function 函数的入参增加到 8 个时, 重写编译程序会得到不同的汇编代码
int my_function(int arg1, int arg2,int arg3,int arg4,int arg5,int arg6,int arg7,int arg8){
	return arg1 + arg2 + arg3 + arg4 + arg5 + arg7 + arg8;
}

int main() {
	int i = my_function(1, 2, 3, 4, 5, 6, 7, 8);
}
