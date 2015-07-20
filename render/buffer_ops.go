// GENERATED CODE
// DO NOT MODIFY
package render

import "unsafe"

type Op byte

const (
	OpInvalid = Op(0x00)
	OpBegin   = Op(0x01)
	OpEnd     = Op(0x02)
	OpColor   = Op(0x03)
	OpPoint   = Op(0x04)
	OpTri     = Op(0x05)
	OpQuad    = Op(0x06)
)

func (op Op) Size() int {
	switch op {
	case OpBegin:
		return int(unsafe.Sizeof(Begin{}))
	case OpEnd:
		return int(unsafe.Sizeof(End{}))
	case OpColor:
		return int(unsafe.Sizeof(Color{}))
	case OpPoint:
		return int(unsafe.Sizeof(Point{}))
	case OpTri:
		return int(unsafe.Sizeof(Tri{}))
	case OpQuad:
		return int(unsafe.Sizeof(Quad{}))
	}
	panic("invalid op")
}

func (w *Buffer) Begin() *Begin { return (*Begin)(w.alloc(OpBegin, int(unsafe.Sizeof(Begin{})))) }
func (r *Reader) Begin() *Begin { return (*Begin)(r.ptr()) }

func (w *Buffer) End() *End { return (*End)(w.alloc(OpEnd, int(unsafe.Sizeof(End{})))) }
func (r *Reader) End() *End { return (*End)(r.ptr()) }

func (w *Buffer) Color() *Color { return (*Color)(w.alloc(OpColor, int(unsafe.Sizeof(Color{})))) }
func (r *Reader) Color() *Color { return (*Color)(r.ptr()) }

func (w *Buffer) Point() *Point { return (*Point)(w.alloc(OpPoint, int(unsafe.Sizeof(Point{})))) }
func (r *Reader) Point() *Point { return (*Point)(r.ptr()) }

func (w *Buffer) Tri() *Tri { return (*Tri)(w.alloc(OpTri, int(unsafe.Sizeof(Tri{})))) }
func (r *Reader) Tri() *Tri { return (*Tri)(r.ptr()) }

func (w *Buffer) Quad() *Quad { return (*Quad)(w.alloc(OpQuad, int(unsafe.Sizeof(Quad{})))) }
func (r *Reader) Quad() *Quad { return (*Quad)(r.ptr()) }
