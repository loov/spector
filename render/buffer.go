package render

import "unsafe"

//go:generate go run buffer_gen.go

type Buffer struct{ data []byte }

func (b *Buffer) alloc(op Op, size int) unsafe.Pointer {
	b.data = append(b.data, byte(op))
	h := len(b.data)
	b.data = append(b.data, make([]byte, size)...)
	return unsafe.Pointer(&b.data[h])
}

func (b *Buffer) Bytes() []byte   { return b.data }
func (b *Buffer) Reader() *Reader { return &Reader{b.data} }

type Reader struct{ data []byte }

func (r *Reader) Op() Op { return Op(r.data[0]) }

func (r *Reader) Next() bool {
	sz := 1 + r.Op().Size()
	r.data = r.data[sz:]
	return len(r.data) > 0
}

func (r *Reader) ptr() unsafe.Pointer {
	return unsafe.Pointer(&r.data[1])
}
