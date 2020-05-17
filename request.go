package httparse

import (
	"bytes"
	"fmt"

	"github.com/pingcap/errors"
)

type Request struct {
	Method  []byte
	Proto   []byte
	URI     []byte
	Headers map[string][][]byte
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
		return 0, errors.New("http version error")
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
		if v, found := h.Headers[b2s(s.key)]; found {
			v = append(v, s.value)
		} else {
			h.Headers[b2s(s.key)] = [][]byte{s.value}
		}
	}
	if s.err != nil {
		return 0, errors.WithStack(s.err)
	}
	return s.hLen, nil
}
