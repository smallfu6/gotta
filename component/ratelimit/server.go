package main

/*

	借助wrk测试接口的QPS, 即基准测试; TODO: wrk
	wrk -c 10 -d 10s -t10 http://localhost:9090

	Running 10s test @ http://localhost:9090
	10 threads and 10 connections
	Thread Stats   Avg      Stdev     Max   +/- Stdev
	Latency   178.42us  589.42us  20.17ms   96.88%
	Req/Sec    10.97k     1.24k   15.13k    76.62%
	1097242 requests in 10.10s, 133.94MB read
	Requests/sec: 108646.15
	Transfer/sec:     13.26MB
*/

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", sayHello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func sayHello(wr http.ResponseWriter, r *http.Request) {
	wr.WriteHeader(200) // TODO
	io.WriteString(wr, "hello world")
}
