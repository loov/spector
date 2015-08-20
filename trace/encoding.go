package trace

type Reader struct {
	Head int
	Data []byte
}

type Writer struct{ Data []byte }

func NewWriter() *Writer { return &Writer{} }

func NewReader() *Reader { return &Reader{} }

func (r *Reader) Append(data []byte) { r.Data = append(r.Data, data...) }

func (r *Reader) readID() ID   { return ID(r.readInt()) }
func (w *Writer) writeID(v ID) { w.writeInt(int32(v)) }

func (r *Reader) readTime() Time   { return Time(r.readInt()) }
func (w *Writer) writeTime(v Time) { w.writeInt(int32(v)) }

func (r *Reader) readFreq() Freq   { return Freq(r.readInt()) }
func (w *Writer) writeFreq(v Freq) { w.writeInt(int32(v)) }

func (r *Reader) readKind() Kind   { return Kind(r.readByte()) }
func (w *Writer) writeKind(v Kind) { w.writeByte(byte(v)) }

func (w *Writer) writeByte(v byte) { w.Data = append(w.Data, v) }
func (r *Reader) readByte() byte {
	v := r.Data[r.Head]
	r.Head++
	return v
}

func (w *Writer) writeInt(v int32) {
	w.Data = append(w.Data,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v>>0),
	)
}
func (r *Reader) readInt() int32 {
	v := int32(r.Data[r.Head+0])<<24 |
		int32(r.Data[r.Head+1])<<16 |
		int32(r.Data[r.Head+2])<<8 |
		int32(r.Data[r.Head+3])<<0
	r.Head += 4
	return v
}

func (w *Writer) writeBytes(v []byte) {
	w.WriteInt(v)
	w.Data = append(w.Data, v...)
}

func (r *Reader) readBytes() []byte {
	sz := r.ReadInt()
	v := r.Data[r.Head : r.Head+sz]
	r.Head += sz
	return v
}

func (w *Writer) writeString(v string) {
	w.WriteBlob([]byte(v))
}

func (r *Reader) readString() string {
	return string(r.ReadBlob())
}
