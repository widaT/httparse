package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/widaT/httparse"
)

type RequestHeader struct {
	*httparse.Request
	noHTTP11        bool
	connectionClose bool
	contentLength   int
	host            []byte
	contentType     []byte
	userAgent       []byte
}

func (r *RequestHeader) Read(buf *bufio.Reader) error {
	b, err := buf.Peek(buf.Buffered())
	if err != nil {
		return errors.WithStack(err)
	}

	n, err := r.Parse(b)
	if err != nil {
		return errors.WithStack(err)
	}
	buf.Discard(n)
	return nil
}
func appendHandle(b []byte, res string) []byte {
	return appendResp(b, "200 OK", "", res)
}

func appendResp(b []byte, status, head, body string) []byte {
	b = append(b, "HTTP/1.1"...)
	b = append(b, ' ')
	b = append(b, status...)
	b = append(b, '\r', '\n')
	b = append(b, "Server: gnet\r\n"...)
	b = append(b, "Date: "...)
	b = time.Now().AppendFormat(b, "Mon, 02 Jan 2006 15:04:05 GMT")
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, "Content-Length: "...)
		b = strconv.AppendInt(b, int64(len(body)), 10)
		b = append(b, '\r', '\n')
	}
	b = append(b, head...)
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, body...)
	}
	return b
}

func main() {

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := l.Accept()
		headle := func(conn net.Conn) {
			fmt.Printf("remote addr %s", conn.RemoteAddr())
			buf := bufio.NewReader(conn)

			for {
				b, err := buf.Peek(4)
				if len(b) == 0 {
					log.Printf("%s", err)
					if err == io.EOF {

						return
					}

					if err == nil {
						panic("bufio.Reader.Peek() returned nil, nil")
					}

					// n == 1 on the first read for the request.
					/* if n == 1 {
						// We didn't read a single byte.
						return errNothingRead{err}
					} */

					return
				}
				rq := RequestHeader{
					Request: &httparse.Request{},
				}

				fmt.Println("---", buf.Buffered())
				b, err = buf.Peek(buf.Buffered())
				if err != nil {
					log.Fatal(err)
				}
				//	fmt.Printf("%s\n", b)

				n, err := rq.Parse(b)
				if err != nil {
					log.Printf("%#v \n", err)
					return
				}
				fmt.Println(n)
				fmt.Printf("%s", rq)

				buf.Discard(n)

				bff := make([]byte, 1024)
				conn.Write(appendHandle(bff, "ddd"))
			}

			//conn.Close()
		}

		go headle(conn)
	}

}
