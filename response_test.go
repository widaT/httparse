package httparse

import (
	"bytes"
	"testing"
)

func TestResponse(t *testing.T) {
	b := []byte("HTTP/1.0 403 Forbidden\nServer: foo.bar\n\n")
	r := Response{}
	n, err := r.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(" %d %s headers len:%d\n", r.StatusCode, r.Proto, len(r.Headers))

	if bytes.Compare(r.Proto, []byte("HTTP/1.0")) != 0 {
		t.Errorf("Proto expect %s got %s", []byte("HTTP/1.0"), r.Proto)
	}

	if r.StatusCode != 403 {
		t.Errorf("StatusCode expect %d got %d", 403, r.StatusCode)
	}

	if n != len(b) {
		t.Errorf("read len expect %d got %d", len(b), n)
	}
}
