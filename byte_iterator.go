package httparse

import (
	"github.com/pingcap/errors"
)

var HeaderNameMap [256]uint8 = [256]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func isHeaderNameToken(b uint8) bool {
	return HeaderNameMap[b] != 0
}

type ByteIter struct {
	data []byte
	pos  int
	len  int
	v    uint8
}

func parseHeaders(b []byte, headers map[string][][]byte, normalizeKey bool) (int, error) {
	iter := NewByteIter(b)
	count := 0
	for iter.next() {
		switch iter.v {
		case '\r':
			if iter.next() {
				if iter.v != '\n' {
					return 0, errors.New("newline error")
				}
				return count + iter.pos, nil
			}
		case '\n':
			return count + iter.pos, nil
		default:
			if !isHeaderNameToken(iter.v) {
				return 0, errors.New("header name error")
			}
			var headerName string
			var value []byte
		key:
			for iter.next() {
				if iter.v == ':' {
					count += iter.pos
					keyData := iter.skip(1)
					if normalizeKey {
						normalizeHeaderKey(keyData)
					}
					headerName = b2s(keyData)
					break key
				} else if !isHeaderNameToken(iter.v) {
					return 0, errors.New("header name error")
				}
			}
		whitespace:
			for iter.next() {
				if iter.v == ' ' || iter.v == '\t' {
					count += iter.pos
					//move
					iter.skip(0)
					continue whitespace
				} else {
					//unread 1
					iter.unRead(1)
					break whitespace
				}
			}

		headerValue:
			for iter.next() {
				switch iter.v {
				case '\r':
					if !iter.next() || iter.v != '\n' {
						return 0, errors.New("got header value newline error")
					}
					count += iter.pos
					value = iter.skip(2)
					break headerValue
				case '\n':
					count += iter.pos
					value = iter.skip(1)
					break headerValue
				}
			}

			if v, found := headers[headerName]; found {
				v = append(v, value)
			} else {
				headers[headerName] = [][]byte{value}
			}
		}

	}
	return count, nil
}

func NewByteIter(b []byte) *ByteIter {
	return &ByteIter{
		data: b,
		len:  len(b),
	}
}

func (iter *ByteIter) next() bool {
	if iter.pos < iter.len {
		iter.v = iter.data[iter.pos]
		iter.pos++
		return true
	}
	return false
}

func (iter *ByteIter) unRead(n int) {
	if iter.pos >= n {
		iter.pos -= n
	}
}

func (iter *ByteIter) skip(skip int) (b []byte) {
	if iter.pos < skip {
		panic("iter skip is bigger than pos")
	}
	headPos := iter.pos - skip
	b = iter.data[:headPos]
	iter.data = iter.data[iter.pos:]
	iter.pos = 0
	iter.len = len(iter.data)
	return
}
