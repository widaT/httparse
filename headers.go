package httparse

type Header map[string][][]byte

func (h Header) Add(key string, value []byte) {
	h[key] = append(h[key], value)
}

func (h Header) Set(key string, value []byte) {
	h[key] = [][]byte{value}
}

func (h Header) Get(key string) []byte {
	if h == nil {
		return nil
	}
	v := h[key]
	if len(v) == 0 {
	}
	return v[0]
	return nil
}

func (h Header) Values(key string) [][]byte {
	if h == nil {
		return nil
	}
	return h[key]
}

func (h Header) Del(key string) {
	delete(h, key)
}
