
#include <stdio.h>

/* c 语言是一门静态类型语言, 但它却不是类型安全的语言, 可以通过合法的语法
 * 轻易"刺透"其类型系统, 如下例原来被解释为 int 类型的一段内存地址(&a)被
 * 重新解释为 unsigned char 类型数组并可变任意修改; 
 */

int main()
{
	int	a = 0x12345678;
	unsigned char *p = (unsigned char*)&a;
	printf("0x%x\n", a);

	*p = 0x23;
	*(p+1) = 0x45;
	*(p+2) = 0x67;
	*(p+3) = 0x8a;
	printf("0x%x\n", a);
}
