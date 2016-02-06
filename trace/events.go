package trace

//go:generate go run events_encoding_gen.go

type Event interface {
	Code() byte
	Decode(dec *Decoder)
	Encode(enc *Encoder)
}

type Handler interface {
	Handle(Event)
}

type Reader interface {
	Next() []Event
}

// Primitive types
type (
	ID   int32
	Time int32
	Freq int32
	Kind byte
)

const (
	InvalidID = ID(0)

	MinTime = Time(-1 << 31)
	MaxTime = Time(1<<31 - 1)
)

func (t Time) Unassigned() bool { return (t == MaxTime) || (t == MinTime) }

const (
	KindInvalid = Kind(0x00)
	KindThread  = Kind(0x01)
	KindStack   = Kind(0x02)

	KindText  = Kind(0x10)
	KindJSON  = Kind(0x11)
	KindBLOB  = Kind(0x12)
	KindImage = Kind(0x13)

	KindUser = Kind(0x20)
)

// Events
type (
	// event: 0x00
	Invalid struct {
	}

	// event: 0x01
	StreamStart struct {
		ProcessID ID
		MachineID ID
		Time      Time
		Freq      Freq
	}

	// event: 0x02
	StreamStop struct {
		Time Time
	}

	// event: 0x03
	ThreadStart struct {
		Time     Time
		ThreadID ID
		StackID  ID
	}

	// event: 0x04
	ThreadSleep struct {
		Time     Time
		ThreadID ID
		StackID  ID
	}

	// event: 0x05
	ThreadWake struct {
		Time     Time
		ThreadID ID
		StackID  ID
	}

	// event: 0x06
	ThreadStop struct {
		Time     Time
		ThreadID ID
		StackID  ID
	}

	// event: 0x07
	SpanBegin struct {
		Time     Time
		ThreadID ID
		StackID  ID
		ID       ID
	}

	// event: 0x08
	SpanEnd struct {
		Time     Time
		ThreadID ID
		StackID  ID
		ID       ID
	}

	// event: 0x09
	Start struct {
		Time     Time
		ThreadID ID
		StackID  ID
		ID       ID
	}

	// event: 0x0A
	Finish struct {
		Time     Time
		ThreadID ID
		StackID  ID
		ID       ID
	}

	// event: 0x0C
	Snapshot struct {
		Time     Time
		ThreadID ID
		StackID  ID
		ID       ID
		Kind     Kind
		Content  []byte
	}

	// event: 0x0D
	Info struct {
		ID      ID
		Name    string
		Kind    Kind
		Content []byte
	}
)
