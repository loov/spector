// GENERATED CODE
// DO NOT MODIFY
package render

import "unsafe"

type Op byte

const (
	OpInvalid = Op(0x00)
	OpColor   = Op(0x01)
	OpRect    = Op(0x02)
	OpTri     = Op(0x03)
	OpQuad    = Op(0x04)
	OpLine    = Op(0x05)
	OpBezier  = Op(0x06)
)

func (op Op) Size() int {
	switch op {
	case OpColor:
		return int(unsafe.Sizeof(Color{}))
	case OpRect:
		return int(unsafe.Sizeof(Rect{}))
	case OpTri:
		return int(unsafe.Sizeof(Tri{}))
	case OpQuad:
		return int(unsafe.Sizeof(Quad{}))
	case OpLine:
		return int(unsafe.Sizeof(Line{}))
	case OpBezier:
		return int(unsafe.Sizeof(Bezier{}))
	}
	panic("invalid op")
}

func (w *Buffer) Color() *Color { return (*Color)(w.alloc(OpColor, int(unsafe.Sizeof(Color{})))) }
func (r *Reader) Color() *Color { return (*Color)(r.ptr()) }

func (w *Buffer) Rect() *Rect { return (*Rect)(w.alloc(OpRect, int(unsafe.Sizeof(Rect{})))) }
func (r *Reader) Rect() *Rect { return (*Rect)(r.ptr()) }

func (w *Buffer) Tri() *Tri { return (*Tri)(w.alloc(OpTri, int(unsafe.Sizeof(Tri{})))) }
func (r *Reader) Tri() *Tri { return (*Tri)(r.ptr()) }

func (w *Buffer) Quad() *Quad { return (*Quad)(w.alloc(OpQuad, int(unsafe.Sizeof(Quad{})))) }
func (r *Reader) Quad() *Quad { return (*Quad)(r.ptr()) }

func (w *Buffer) Line() *Line { return (*Line)(w.alloc(OpLine, int(unsafe.Sizeof(Line{})))) }
func (r *Reader) Line() *Line { return (*Line)(r.ptr()) }

func (w *Buffer) Bezier() *Bezier { return (*Bezier)(w.alloc(OpBezier, int(unsafe.Sizeof(Bezier{})))) }
func (r *Reader) Bezier() *Bezier { return (*Bezier)(r.ptr()) }
