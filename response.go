package httparse

import (
	"fmt"

	"github.com/pkg/errors"
)

type Response struct {
	Proto      []byte
	StatusCode int
	//Reason  []byte
	Headers            map[string][][]byte
	normalizeHeaderKey bool
}

func (r *Response) Reset() {
	r.Proto = nil
	r.StatusCode = 0
	r.Headers = nil
}

//for debug
func (r Response) String() string {
	str := fmt.Sprintf("StatusCode:%d,Proto:%s\n", r.StatusCode, r.Proto)
	str += "Headers\n"
	for k, v := range r.Headers {
		str += fmt.Sprintf("  %s:%s \n", k, v)
	}
	return str
}

func (h *Response) Parse(b []byte) (int, error) {
	if len(b) < 12 {
		return 0, errors.New("to short")
	}

	length := 0
	//@todo skip enpty lines
	if err := checkVersion(b); err != nil {
		return 0, errors.New("checkVersion err")
	}

	h.Proto = b[:8]
	if b[8] != ' ' {
		return 0, errors.New("parse error")
	}

	//parse code
	a := func(n uint8) (int, error) {
		n -= '0'
		if n > 9 {
			return 0, errors.New("status code err Expecting 0-9")
		}
		return int(n), nil
	}
	hundreds, err := a(b[9])
	if err != nil {
		return 0, err
	}
	tens, err := a(b[10])
	if err != nil {
		return 0, err
	}
	ones, err := a(b[11])
	if err != nil {
		return 0, err
	}
	h.StatusCode = hundreds*100 + tens*10 + ones

	//skip reason, read to the end of first line
	length += 12
	for i := 12; i < len(b); i++ {
		length++
		if b[i] == '\n' {
			break
		}
	}
	n, err := h.parseHeaders(b[length:])
	if err != nil {
		return 0, err
	}
	return length + n, nil
}

func (h *Response) parseHeaders(buf []byte) (int, error) {
	if h.Headers == nil {
		h.Headers = make(map[string][][]byte)
	}
	return parseHeaders(buf, h.Headers, h.normalizeHeaderKey)
}
