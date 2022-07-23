package main

import (
	"log"
	"net"
	"time"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}

	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var sl []net.Conn
	for i := 0; i < 1000; i++ {
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}

	time.Sleep(time.Second * 120)
}
