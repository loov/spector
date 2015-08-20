package trace

type Decoder struct {
	Head int
	Data []byte
}

type Encoder struct{ Data []byte }

func NewEncoder() *Encoder { return &Encoder{} }
func NewDecoder() *Decoder { return &Decoder{} }

func (dec *Decoder) Append(data []byte) { dec.Data = append(dec.Data, data...) }

func (dec *Decoder) readID() ID   { return ID(dec.readInt()) }
func (enc *Encoder) writeID(v ID) { enc.writeInt(int32(v)) }

func (dec *Decoder) readTime() Time   { return Time(dec.readInt()) }
func (enc *Encoder) writeTime(v Time) { enc.writeInt(int32(v)) }

func (dec *Decoder) readFreq() Freq   { return Freq(dec.readInt()) }
func (enc *Encoder) writeFreq(v Freq) { enc.writeInt(int32(v)) }

func (dec *Decoder) readKind() Kind   { return Kind(dec.readByte()) }
func (enc *Encoder) writeKind(v Kind) { enc.writeByte(byte(v)) }

func (enc *Encoder) writeByte(v byte) { enc.Data = append(enc.Data, v) }
func (dec *Decoder) readByte() byte {
	v := dec.Data[dec.Head]
	dec.Head++
	return v
}

func (enc *Encoder) writeInt(v int32) {
	enc.Data = append(enc.Data, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
func (dec *Decoder) readInt() int32 {
	d, h := dec.Data, dec.Head
	v := int32(d[h])<<24 | int32(d[h+1])<<16 | int32(d[h+2])<<8 | int32(d[h+3])
	dec.Head += 4
	return v
}

func (enc *Encoder) writeBytes(v []byte) {
	enc.writeInt(int32(len(v)))
	enc.Data = append(enc.Data, v...)
}
func (dec *Decoder) readBytes() []byte {
	sz := int(dec.readInt())
	v := dec.Data[dec.Head : dec.Head+sz]
	dec.Head += sz
	return v
}

func (enc *Encoder) writeString(v string) { enc.writeBytes([]byte(v)) }
func (dec *Decoder) readString() string   { return string(dec.readBytes()) }
