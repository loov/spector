package render

import "unsafe"

//go:generate go run buffer_gen.go

type Buffer struct{ data []byte }

func NewBuffer(size int) *Buffer {
	return &Buffer{make([]byte, 0, size)}
}

func (b *Buffer) Reset() { b.data = b.data[:0:cap(b.data)] }

func (b *Buffer) alloc(op Op, size int) unsafe.Pointer {
	b.data = append(b.data, byte(op))
	h := len(b.data)
	if len(b.data)+size+1 < cap(b.data) {
		b.data = b.data[:len(b.data)+size]
	} else {
		b.data = append(b.data, make([]byte, size)...)
	}
	return unsafe.Pointer(&b.data[h])
}

func (b *Buffer) Len() int        { return len(b.data) }
func (b *Buffer) Bytes() []byte   { return b.data }
func (b *Buffer) Reader() *Reader { return &Reader{0, b.data} }

type Reader struct {
	head int
	data []byte
}

func (r *Reader) Op() Op { return Op(r.data[r.head]) }

func (r *Reader) Next() bool {
	sz := 1 + r.Op().Size()
	r.head += sz
	return r.head < len(r.data)
}

func (r *Reader) ptr() unsafe.Pointer {
	return unsafe.Pointer(&r.data[r.head+1])
}
