package httparse

import (
	"bytes"
	"testing"
)

func TestRequset(t *testing.T) {
	b := []byte("GET / HTTP/1.1\r\n\r\n")

	t.Log(len(b))

	r := Request{}
	n, err := r.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	if bytes.Compare(r.Method, []byte("GET")) != 0 {
		t.Errorf("metcho expect %s got %s", []byte("GET"), r.Method)
	}

	if bytes.Compare(r.Proto, []byte("HTTP/1.1")) != 0 {
		t.Errorf("proto expect %s got %s", []byte("HTTP/1.1"), r.Proto)
	}
	if bytes.Compare(r.URI, []byte("/")) != 0 {
		t.Errorf("uri expect %s got %s", []byte("/"), r.URI)
	}

	if n != len(b) {
		t.Errorf("read len expect %d got %d", len(b), n)
	}
}
