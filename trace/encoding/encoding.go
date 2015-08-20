package encoding

type Reader struct {
	Head int
	Data []byte
}

type Writer struct {
	Data []byte
}

func NewReader(data []byte) *Reader { return &Reader{0, data} }
func NewWriter() *Writer            { return &Writer{} }

func (w *Writer) WriteByte(v byte) {
	w.Data = append(w.Data, v)
}

func (r *Reader) ReadByte() byte {
	v := r.Data[r.Head]
	r.Head++
	return v
}

func (w *Writer) WriteInt(v int32) {
	w.Data = append(w.Data,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v>>0),
	)
}

func (r *Reader) ReadInt() int32 {
	v := int32(r.Data[r.Head+0])<<24 |
		int32(r.Data[r.Head+1])<<16 |
		int32(r.Data[r.Head+2])<<8 |
		int32(r.Data[r.Head+3])<<0
	r.Head += 4
	return v
}

func (w *Writer) WriteBlob(v []byte) {
	w.WriteInt(v)
	w.Data = append(w.Data, v...)
}

func (r *Reader) ReadBlob() []byte {
	sz := r.ReadInt()
	v := r.Data[r.Head : r.Head+sz]
	r.Head += sz
	return v
}

func (w *Writer) WriteUTF8(v string) {
	w.WriteBlob([]byte(v))
}

func (r *Reader) ReadUTF8() string {
	return string(r.ReadBlob())
}
