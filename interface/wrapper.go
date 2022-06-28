package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

/*
	包裹函数(wrapper function)
	接受接口类型参数, 并返回与其参数类型相同的返回值;
	TODO: 对比函数式编程概念函子的用法

	通过包裹函数可以实现对输入数据的过滤、装饰、变换等操作,
	并将结果再次返回给调用者;
	由于包裹函数的返回值类型与参数类型相同, 因此我们可以将多个接受同一接口
	类型参数的包裹函数组合成一条链来调用;
*/

func CapReader(r io.Reader) io.Reader {
	return &capitaliizedReader{r: r}
}

type capitaliizedReader struct {
	r io.Reader
}

func (r *capitaliizedReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err != nil {
		return 0, err
	}

	q := bytes.ToUpper(p)
	for i, v := range q {
		p[i] = v
	}
	return n, err
}

func main() {
	r := strings.NewReader("hello, gopher!\n")
	r1 := CapReader(io.LimitReader(r, 4))
	if _, err := io.Copy(os.Stdout, r1); err != nil {
		log.Fatal(err)
	}
	// 将CapReader和io.LimitReader串在一起形成了一条调用链, 这条调用链的
	// 功能为: 截取输入数据的前4字节并将其转换为大写字母;

}
