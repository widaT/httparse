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

func TestParse(t *testing.T) {
	b := []byte("GET / HTTP/1.1\r\nHost: foo.com\r\nCookie: \r\n\r\n")
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

	if len(r.Headers) != 2 {
		t.Errorf("Headers expect 2 got %d", len(r.Headers))
	}

	v := r.Headers.Get("Host")
	if len(v) == 0 || !bytes.Equal(v, []byte("foo.com")) {
		t.Errorf("read host err %s -- %#v", v, r.Headers)
	}

	v = r.Headers.Get("Cookie")
	if !bytes.Equal(v, []byte("")) {
		t.Errorf("read cookie err %d", v)
	}
}

func TestToUpperTable(t *testing.T) {

	s := "AcdDSS"

	t.Errorf("----%s\n", []byte{toLowerTable(s[4])})
	t.Errorf("----%s\n", []byte{toUpperTable(s[0])})
}
