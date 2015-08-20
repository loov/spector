package trace

import "github.com/egonelbre/spector/trace/encoding"

type Reader struct{ enc encoding.Reader }

func NewReader() *Reader {
	return &Reader{*encoding.NewReader([]byte{})}
}

func (r *Reader) Append(data []byte) { r.enc.Data = append(r.enc.Data, data...) }
