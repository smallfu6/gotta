package main

import (
	"fmt"
	"net/http"
)

/*
	TODO: tls, ssl 基础知识, openssl 命令的掌握
	使用 http.ListenAndServeTLS 函数实现http的安全通信
*/

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})
	fmt.Println(http.ListenAndServeTLS("localhost:8080", "server.crt", "server.key", nil))
}

/*
	生成私钥
	openssl genrsa -out server.key 2048

	生成公钥
	openssl req -new -x509 -key server.key  -out server.crt -days 365


	curl https://localhost:8080
	curl: (60) SSL certificate problem: self signed certificate
	报错的主要原因是示例中HTTPS Web服务所使用证书(server.crt)是我们自己生成
	的自签名证书, curl使用测试环境系统中内置的各种数字证书授权机构的公钥证书
	无法对其进行验证;
	curl -k https://localhost:8080
	使用-k选项忽略对示例中HTTPS Web服务的服务端证书的校验
*/
