package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
	中间件(middleware)
	含义可大可小, 在Go Web编程中常常指的是一个实现了http.Handler接口的
	http.HandlerFunc类型实例, 实质上, 这里的中间件就是包裹函数和适配器
	函数类型结合的产物;
*/

func validateAuth(s string) error {
	if s != "123456" {
		return fmt.Errorf("%s", "bad auth token")
	}
	return nil
}

func greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t)
		h.ServeHTTP(w, r)
	})
}

func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validateAuth(r.Header.Get("auth"))
		if err != nil {
			http.Error(w, "bad auth param", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})

}

func main() {
	http.ListenAndServe(":8080", logHandler(authHandler(http.HandlerFunc(greetings))))
	// logHandler、authHandler 本质上就是一个包裹函数(支持链式调用), 但其内
	// 部利用了适配器函数类型(http.HandlerFunc)将一个普通函数(如例子中的几个
	// 匿名函数)转换为实现了http.Handler的类型的实例, 并将其作为返回值返回;
}
