package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func establishConn() net.Conn {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Printf("dial error: %s", err)
		return nil
	}

	log.Println(":connect to server ok")
	return conn
}

func main() {
	conn := establishConn()
	if conn == nil {
		return
	}
	defer conn.Close()
	var i int
	for {
		if i < 10 {
			message := fmt.Sprintf("test %d", i)
			conn.Write([]byte(message))
			log.Println("send data: ", message)
		}
		time.Sleep(1 * time.Second)
		i++
	}
}
