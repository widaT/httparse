package httparse

import (
	"bytes"
	"fmt"

	"github.com/pingcap/errors"
)

type Request struct {
	Method             []byte
	Proto              []byte
	URI                []byte
	Headers            Header
	normalizeHeaderKey bool
}

func NewRequst() *Request {
	return &Request{
		normalizeHeaderKey: true,
	}
}

func (r *Request) Reset() {
	r.Headers = nil
	r.normalizeHeaderKey = true
}

//for debug
func (r Request) String() string {
	str := fmt.Sprintf("Method:%s,Proto:%s,URI:%s\n", r.Method, r.Proto, r.URI)
	str += "Headers\n"
	for k, v := range r.Headers {
		str += fmt.Sprintf("  %s:%s \n", k, v)
	}
	return str
}

func (r *Request) NormalizeHeaderKey(b bool) {
	r.normalizeHeaderKey = b
}

func (h *Request) Parse(b []byte) (int, error) {
	var length int = 0
	//@todo skip empty lines

	//parse method
	for i := 0; i < len(b); i++ {
		length++
		if b[i] == ' ' {
			h.Method = b[:i]
			b = b[i+1:]
			break
		} else if !isToken(b[i]) {
			return 0, errors.New("token error")
		}
	}
	//parse uri
	for i := 0; i < len(b); i++ {
		length++
		if b[i] == ' ' {
			h.URI = b[:i]
			b = b[i+1:]
			break
		}
	}

	//parse httpversion
	if len(b) < 10 {
		return 0, errors.Errorf("parse http version error want len >10 got %d", len(b))
	}

	if err := checkVersion(b); err != nil {
		return 0, err
	}
	h.Proto = b[:8]
	//newline
	if b[8] != '\r' || b[9] != '\n' {
		return 0, errors.New("http version error")
	}
	length += 10
	b = b[10:]
	len, err := h.parseHeaders(b)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return length + len, nil
}

func (r *Request) DelHeader(key string) {
	r.Headers.Del(key)
}

func checkVersion(b []byte) error {
	if len(b) < 8 {
		return errors.New("too short")
	}
	if bytes.Compare(b[:7], []byte("HTTP/1.")) != 0 {
		return errors.New("unsport http version")
	}
	if b[7] != '0' && b[7] != '1' {
		return errors.New("http version error")
	}
	return nil
}

func isToken(b uint8) bool {
	return b > 0x1F && b < 0x7F
}

func (h *Request) parseHeaders(buf []byte) (int, error) {
	if h.Headers == nil {
		h.Headers = make(Header)
	}
	return parseHeaders(buf, h.Headers, h.normalizeHeaderKey)
}

func (h *Request) GetHeader(key string) []byte {
	return h.Headers.Get(key)
}

func (h *Request) SetHeader(key string, val []byte) {
	if h.Headers == nil {
		h.Headers = make(Header)
	}
	h.Headers.Set(key, val)
}

func (h *Request) AddHeader(key string, val []byte) {
	if h.Headers == nil {
		h.Headers = make(Header)
	}
	h.Headers.Add(key, val)
}

func normalizeHeaderKey(b []byte) {
	n := len(b)
	if n == 0 {
		return
	}
	b[0] = toUpperTable(b[0])
	for i := 1; i < n; i++ {
		p := &b[i]
		if *p == '-' {
			i++
			if i < n {
				b[i] = toUpperTable(b[i])
			}
			continue
		}
		*p = toLowerTable(*p)
	}
}

func toUpperTable(letter uint8) (r uint8) {
	r = letter
	if letter >= 'a' && letter <= 'z' {
		r -= 32
	}
	return
}

func toLowerTable(letter uint8) (r uint8) {
	r = letter
	if 'A' <= letter && letter <= 'Z' {
		r += 32
	}
	return
}
