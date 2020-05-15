package main

import (
	"fmt"
	"log"
	"net"

	"github.com/widaT/httparse"
)

func main() {

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := l.Accept()
		headle := func(conn net.Conn) {
			buf := make([]byte, 30*1024)
			n, _ := conn.Read(buf)
			r := httparse.Request{}
			num, err := r.Parse(buf[:n])
			if err != nil {
				log.Fatalf("%#v \n", err)
				return
			}
			fmt.Println(n, num)
			fmt.Printf("%s", r)

			conn.Write([]byte("HTTP/1.1 \r\n 200 Ok\r\n\r\n"))
			conn.Close()
		}

		go headle(conn)
	}

}
