package httparse

import (
	"bytes"
	"errors"
	"unsafe"
)

var (
	errNeedMore       = errors.New("need more data: cannot find trailing lf")
	errInvalidName    = errors.New("invalid header name")
	errEmptyHexNum    = errors.New("empty hex number")
	errTooLargeHexNum = errors.New("too large hex number")
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//code from  fasthttp
type headerScanner struct {
	b     []byte
	key   []byte
	value []byte
	err   error

	hLen        int
	nextColon   int
	nextNewLine int

	initialized bool
}

//code from fasthttp
func (s *headerScanner) next() bool {
	if !s.initialized {
		s.nextColon = -1
		s.nextNewLine = -1
		s.initialized = true
	}
	bLen := len(s.b)
	if bLen >= 2 && s.b[0] == '\r' && s.b[1] == '\n' {
		s.b = s.b[2:]
		s.hLen += 2
		return false
	}
	if bLen >= 1 && s.b[0] == '\n' {
		s.b = s.b[1:]
		s.hLen++
		return false
	}
	var n int
	if s.nextColon >= 0 {
		n = s.nextColon
		s.nextColon = -1
	} else {
		n = bytes.IndexByte(s.b, ':')

		// There can't be a \n inside the header name, check for this.
		x := bytes.IndexByte(s.b, '\n')
		if x < 0 {
			// A header name should always at some point be followed by a \n
			// even if it's the one that terminates the header block.
			s.err = errNeedMore
			return false
		}
		if x < n {
			// There was a \n before the :
			s.err = errInvalidName
			return false
		}
	}
	if n < 0 {
		s.err = errNeedMore
		return false
	}
	s.key = s.b[:n]
	n++
	for len(s.b) > n && s.b[n] == ' ' {
		n++
		// the newline index is a relative index, and lines below trimed `s.b` by `n`,
		// so the relative newline index also shifted forward. it's safe to decrease
		// to a minus value, it means it's invalid, and will find the newline again.
		s.nextNewLine--
	}
	s.hLen += n
	s.b = s.b[n:]
	if s.nextNewLine >= 0 {
		n = s.nextNewLine
		s.nextNewLine = -1
	} else {
		n = bytes.IndexByte(s.b, '\n')
	}
	if n < 0 {
		s.err = errNeedMore
		return false
	}
	for {
		if n+1 >= len(s.b) {
			break
		}
		if s.b[n+1] != ' ' && s.b[n+1] != '\t' {
			break
		}
		d := bytes.IndexByte(s.b[n+1:], '\n')
		if d <= 0 {
			break
		} else if d == 1 && s.b[n+1] == '\r' {
			break
		}
		e := n + d + 1
		if c := bytes.IndexByte(s.b[n+1:e], ':'); c >= 0 {
			s.nextColon = c
			s.nextNewLine = d - c - 1
			break
		}
		n = e
	}
	if n >= len(s.b) {
		s.err = errNeedMore
		return false
	}
	s.value = s.b[:n]
	s.hLen += n + 1
	s.b = s.b[n+1:]

	if n > 0 && s.value[n-1] == '\r' {
		n--
	}
	for n > 0 && s.value[n-1] == ' ' {
		n--
	}
	s.value = s.value[:n]
	return true
}
