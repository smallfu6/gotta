package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"net/http/httptest"
	"time"
)

type result struct {
	value string
}

func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))
	queryFunc := func(server *httptest.Server) {
		defer server.Close()
		url := server.URL
		fmt.Println(url)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("http get error: %s\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
	}
	for _, serv := range servers {
		go queryFunc(serv)
	}
	return <-c, nil
}

// 使用 httptest 包的 NewServer 函数创建了三个模拟的服务端, 然后将这三个服务端
// 的实例传入 first 函数; TODO: httptest 包的使用
func fakeWeatherServer(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s receive a http request\n", name)
			time.Sleep(1 * time.Second)
			w.Write([]byte(name + ":ok"))
		},
	))
}

// func main() {
// 	result, err := first1(
// 		fakeWeatherServer("open-weather-1"),
// 		fakeWeatherServer("open-weather-2"),
// 		fakeWeatherServer("open-weather-3"),
// 	)

// 	if err != nil {
// 		log.Println("invoke first error:", err)
// 		return
// 	}
// 	log.Println(result)
// }

// 在 first 的基础上增加了定时器, 通过 select 原语监视定时器事件
// 和响应 channel 上的事件; 如果响应channel上长时间没有数据返回,
// 则当定时器事件触发时, first1 函数返回超时错误
func first1(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))
	queryFunc := func(server *httptest.Server) {
		defer server.Close()
		url := server.URL
		fmt.Println(url)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("http get error: %s\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
	}
	for _, serv := range servers {
		go queryFunc(serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("timeout")
	}
}

// 加了超时模式的版本依然存在问题, 即使 first1 函数因为超时返回,
// 但是已经创建的 goroutine 可能仍然在请求服务端或等待应答状态,
// 没有返回, 也没有被回收, 资源仍然在占用, 即使它们的存在没有
// 任何意义; 可以使用 context 包使 goroutine 支持取消操作
func first2(servers ...*httptest.Server) (result, error) {
	c := make(chan result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryFunc := func(i int, server *httptest.Server) {
		defer server.Close()
		url := server.URL
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("query goroutine-%d: http NewRequest error: %s\n", i, err)
			return
		}

		req = req.WithContext(ctx)
		// http 包支持利用 context.Context 的超时和 cancel 机制
		// TODO: 熟悉

		log.Printf("query goroutine-%d: send response...\n", i)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("query goroutine-%d: get return error: %s\n", i, err)
			return
		}

		log.Printf("query goroutine-%d: get response\n", i)
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
		return
	}

	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("timeout")
	}
}

func fakeWeatherServer2(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s receive a http request\n", name)
			time.Sleep(time.Duration(interval) * time.Millisecond)
			w.Write([]byte(name + ":ok"))
		},
	))
}

// 利用context.WithCancel创建了一个可以取消的context.Context变量, 在每个
// 发起查询请求的goroutine中, 用该变量更新了request中的ctx变量, 使其支持
// 被取消; 这样在first2函数中, 无论是成功得到某个查询goroutine的返回结果,
// 还是超时失败返回, 通过defer cancel()设定cancel函数在first2函数返回前
// 被执行, 那些尚未返回的在途(on-flight)查询的goroutine都将收到cancel事
// 件并退出;
func main() {
	result, err := first2(
		fakeWeatherServer2("open-weather-1", 200),
		fakeWeatherServer2("open-weather-2", 1000),
		fakeWeatherServer2("open-weather-3", 600),
	)
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	fmt.Println(result)
	time.Sleep(10 * time.Second)
}
