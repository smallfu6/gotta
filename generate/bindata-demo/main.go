package main

import (
	"fmt"
	"net/http"
)

/*
	go generate驱动从静态资源文件数据到Go源码的转换

	go generate结合go-bindata工具(https://github.com/go-bindata/go-bindata)
	将静态资源文件也嵌入可执行文件中, 尤其是在Web开发领域, Gopher希望将一些
	静态资源文件(比如CSS文件等)嵌入最终的二进制文件中一起发布和部署;
	go generate main.go 会生成 ./static.go
	go build main.go static.go

	即使删除static/16130308.jpeg 也不会影响到程序的应答返回结果, 因为图片
	数据已经嵌入二进制程序当中了, 16130308.jpeg 将随着二进制程序一并分发与部署;

*/

//go:generate go-bindata -o  static.go static/16130308.jpeg

func main() {
	data, err := Asset("static/16130308.jpeg")
	if err != nil {
		fmt.Println("Asset invoke error:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write(data)
	})

	fmt.Println("listen server :8080")
	http.ListenAndServe(":8080", nil)
}
