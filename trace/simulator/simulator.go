package simulator

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/egonelbre/spector/trace"
)

type Limits struct {
	Open    int
	Async   int
	Threads int
}

type process struct{ killed int32 }

func (p *process) Alive() bool { return atomic.LoadInt32(&p.killed) == 0 }
func (p *process) Stop()       { atomic.StoreInt32(&p.killed, 1) }

type Stream struct {
	// should not be changed after starting
	max            Limits
	eventInterval  time.Duration
	threadInterval time.Duration

	process

	// owned by the routine
	time         int32
	lastThreadID trace.ID
	threads      []*Thread

	// buffer is protected by mutex
	mu     sync.Mutex
	buffer []trace.Event
}

func NewStream() *Stream {
	stream := &Stream{}
	stream.max = Limits{
		Open:    8,
		Async:   2,
		Threads: 20,
	}
	stream.eventInterval = 100 * time.Millisecond
	stream.threadInterval = 150 * time.Millisecond

	return stream
}

func (stream *Stream) emit(event trace.Event) {
	stream.mu.Lock()
	stream.buffer = append(stream.buffer, event)
	stream.mu.Unlock()
}

func (stream *Stream) Next() (pending []trace.Event) {
	stream.mu.Lock()
	pending, stream.buffer = stream.buffer, nil
	stream.mu.Unlock()
	return
}

func (stream *Stream) Start() { go stream.Run() }

func (stream *Stream) advanceTime(amount int32) trace.Time {
	return trace.Time(atomic.AddInt32(&stream.time, amount))
}

func (stream *Stream) Run() {
	stream.emit(&trace.StreamStart{
		ProcessID: trace.ID(rand.Int31()),
		MachineID: trace.ID(rand.Int31()),
		Time:      stream.advanceTime(0),
		Freq:      1e9, // ticks per second
	})

	for stream.Alive() {
		time.Sleep(stream.threadInterval)
		stream.advanceTime(rand.Int31n(100))

		n := rand.Intn(stream.max.Threads)
		if n < len(stream.threads) {
			thread := stream.threads[n]
			thread.Stop()
			stream.threads = append(stream.threads[:n], stream.threads[n+1:]...)
		} else {
			thread := &Thread{}
			thread.stream = stream
			thread.id = stream.lastThreadID
			stream.lastThreadID++
			stream.threads = append(stream.threads, thread)
			thread.Start()
		}
	}

	for _, thread := range stream.threads {
		thread.Stop()
	}
}

type Thread struct {
	stream *Stream
	process
	id          trace.ID
	lastEventId trace.ID
	open        []trace.ID
}

func (thread *Thread) Start() { go thread.Run() }

func (thread *Thread) Run() {
	stream := thread.stream

	stream.emit(&trace.ThreadStart{
		Time:     stream.advanceTime(rand.Int31n(5)),
		ThreadID: thread.id,
		StackID:  0,
	})

	for thread.Alive() {
		time.Sleep(stream.eventInterval)
		time := stream.advanceTime(rand.Int31n(100))

		n := rand.Intn(stream.max.Open)
		if n < len(thread.open) {
			id := thread.open[n]
			stream.emit(&trace.SpanEnd{
				Time:     time,
				ThreadID: thread.id,
				StackID:  0,
				ID:       id,
			})

			thread.open = append(thread.open[:n], thread.open[n+1:]...)
		} else {
			id := thread.lastEventId
			thread.lastEventId++

			stream.emit(&trace.SpanBegin{
				Time:     time,
				ThreadID: thread.id,
				StackID:  0,
				ID:       id,
			})
			thread.open = append(thread.open, id)
		}
	}

	stream.emit(&trace.ThreadStop{
		Time:     stream.advanceTime(rand.Int31n(5)),
		ThreadID: thread.id,
		StackID:  0,
	})
}
