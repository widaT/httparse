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
	Headers            map[string][][]byte
	normalizeHeaderKey bool
}

func NewRequst() *Request {
	return &Request{
		normalizeHeaderKey: true,
	}
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
	delete(r.Headers, key)
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
	var s headerScanner
	s.b = buf
	for s.next() {
		if h.Headers == nil {
			h.Headers = make(map[string][][]byte)
		}
		if h.normalizeHeaderKey {
			normalizeHeaderKey(s.key)
		}
		key := b2s(s.key)
		if v, found := h.Headers[key]; found {
			v = append(v, s.value)
		} else {
			h.Headers[key] = [][]byte{s.value}
		}
	}
	if s.err != nil {
		return 0, errors.WithStack(s.err)
	}
	return s.hLen, nil
}

func (h *Request) GetHeader(key string) [][]byte {
	if h.Headers == nil {
		return nil
	}
	if v, found := h.Headers[key]; found {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}

const toLowerTable = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@abcdefghijklmnopqrstuvwxyz[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
const toUpperTable = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`ABCDEFGHIJKLMNOPQRSTUVWXYZ{|}~\u007f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"

func normalizeHeaderKey(b []byte) {
	n := len(b)
	if n == 0 {
		return
	}
	b[0] = toUpperTable[b[0]]
	for i := 1; i < n; i++ {
		p := &b[i]
		if *p == '-' {
			i++
			if i < n {
				b[i] = toUpperTable[b[i]]
			}
			continue
		}
		*p = toLowerTable[*p]
	}
}
