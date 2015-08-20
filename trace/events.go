package trace

type Event interface {
	Code() byte
	ReadFrom(r *Reader)
	WriteTo(w *Writer)
}

func NewEventByCode(code byte) Event {
	switch code {
	case 0x00:
		return &Invalid{}
	case 0x01:
		return &StreamStart{}
	case 0x02:
		return &StreamStop{}
	case 0x03:
		return &ThreadStart{}
	case 0x04:
		return &ThreadSleep{}
	case 0x05:
		return &ThreadWake{}
	case 0x06:
		return &ThreadStop{}
	case 0x07:
		return &Begin{}
	case 0x08:
		return &End{}
	case 0x09:
		return &Start{}
	case 0x0A:
		return &Finish{}
	case 0x0C:
		return &Snapshot{}
	case 0x0D:
		return &Info{}
	}
	panic("unknown code")
}

// code: 0x00
type Invalid struct {
}

// code: 0x01
type StreamStart struct {
	ProcessID    int32
	MachineID    int32
	Time         int32
	CPUFrequency int32
}

// code: 0x02
type StreamStop struct {
	Time int32
}

// code: 0x03
type ThreadStart struct {
	Time     int32
	ThreadID int32
	StackID  int32
}

// code: 0x04
type ThreadSleep struct {
	Time     int32
	ThreadID int32
	StackID  int32
}

// code: 0x05
type ThreadWake struct {
	Time     int32
	ThreadID int32
	StackID  int32
}

// code: 0x06
type ThreadStop struct {
	Time     int32
	ThreadID int32
	StackID  int32
}

// code: 0x07
type Begin struct {
	Time     int32
	ThreadID int32
	StackID  int32
	ID       int32
}

// code: 0x08
type End struct {
	Time     int32
	ThreadID int32
	StackID  int32
	ID       int32
}

// code: 0x09
type Start struct {
	Time     int32
	ThreadID int32
	StackID  int32
	ID       int32
}

// code: 0x0A
type Finish struct {
	Time     int32
	ThreadID int32
	StackID  int32
	ID       int32
}

// code: 0x0C
type Snapshot struct {
	Time        int32
	ThreadID    int32
	StackID     int32
	ID          int32
	ContentKind byte
	Content     []byte
}

// code: 0x0D
type Info struct {
	ID          int32
	Name        string
	ContentKind byte
	Content     []byte
}

func (ev *Invalid) Code() byte { return 0x00 }
func (ev *Invalid) ReadFrom(r *Reader) {
}
func (ev *Invalid) WriteTo(w *Writer) {
}
func (ev *StreamStart) Code() byte { return 0x01 }
func (ev *StreamStart) ReadFrom(r *Reader) {
	ev.ProcessID = r.enc.ReadInt()
	ev.MachineID = r.enc.ReadInt()
	ev.Time = r.enc.ReadInt()
	ev.CPUFrequency = r.enc.ReadInt()
}
func (ev *StreamStart) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.ProcessID)
	w.enc.WriteInt(ev.MachineID)
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.CPUFrequency)
}
func (ev *StreamStop) Code() byte { return 0x02 }
func (ev *StreamStop) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
}
func (ev *StreamStop) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
}
func (ev *ThreadStart) Code() byte { return 0x03 }
func (ev *ThreadStart) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
}
func (ev *ThreadStart) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
}
func (ev *ThreadSleep) Code() byte { return 0x04 }
func (ev *ThreadSleep) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
}
func (ev *ThreadSleep) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
}
func (ev *ThreadWake) Code() byte { return 0x05 }
func (ev *ThreadWake) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
}
func (ev *ThreadWake) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
}
func (ev *ThreadStop) Code() byte { return 0x06 }
func (ev *ThreadStop) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
}
func (ev *ThreadStop) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
}
func (ev *Begin) Code() byte { return 0x07 }
func (ev *Begin) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
	ev.ID = r.enc.ReadInt()
}
func (ev *Begin) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
	w.enc.WriteInt(ev.ID)
}
func (ev *End) Code() byte { return 0x08 }
func (ev *End) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
	ev.ID = r.enc.ReadInt()
}
func (ev *End) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
	w.enc.WriteInt(ev.ID)
}
func (ev *Start) Code() byte { return 0x09 }
func (ev *Start) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
	ev.ID = r.enc.ReadInt()
}
func (ev *Start) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
	w.enc.WriteInt(ev.ID)
}
func (ev *Finish) Code() byte { return 0x0A }
func (ev *Finish) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
	ev.ID = r.enc.ReadInt()
}
func (ev *Finish) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
	w.enc.WriteInt(ev.ID)
}
func (ev *Snapshot) Code() byte { return 0x0C }
func (ev *Snapshot) ReadFrom(r *Reader) {
	ev.Time = r.enc.ReadInt()
	ev.ThreadID = r.enc.ReadInt()
	ev.StackID = r.enc.ReadInt()
	ev.ID = r.enc.ReadInt()
	ev.ContentKind = r.enc.ReadByte()
	ev.Content = r.enc.ReadBlob()
}
func (ev *Snapshot) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.Time)
	w.enc.WriteInt(ev.ThreadID)
	w.enc.WriteInt(ev.StackID)
	w.enc.WriteInt(ev.ID)
	w.enc.WriteByte(ev.ContentKind)
	w.enc.WriteBlob(ev.Content)
}
func (ev *Info) Code() byte { return 0x0D }
func (ev *Info) ReadFrom(r *Reader) {
	ev.ID = r.enc.ReadInt()
	ev.Name = r.enc.ReadUTF8()
	ev.ContentKind = r.enc.ReadByte()
	ev.Content = r.enc.ReadBlob()
}
func (ev *Info) WriteTo(w *Writer) {
	w.enc.WriteInt(ev.ID)
	w.enc.WriteUTF8(ev.Name)
	w.enc.WriteByte(ev.ContentKind)
	w.enc.WriteBlob(ev.Content)
}
