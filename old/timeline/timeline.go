package timeline

import (
	"github.com/egonelbre/spector/trace"
)

type Timeline struct {
	Procs       []*Proc
	TotalEvents int64
}

func (timeline *Timeline) ProcByID(pid trace.ID) (*Proc, bool) {
	for _, proc := range timeline.Procs {
		if proc.PID == pid {
			return proc, true
		}
	}
	return nil, false
}

type Proc struct {
	PID  trace.ID
	MID  trace.ID
	Time trace.Time
	Freq trace.Freq

	Start trace.Time
	Stop  trace.Time

	Threads []*Thread
	Tracks  []*Track
}

// OpenTrack gets an open track starting from time
func (proc *Proc) OpenTrack(time trace.Time) *Track {
	for _, track := range proc.Tracks {
		if track.End() <= time {
			return track
		}
	}

	track := &Track{}
	proc.Tracks = append(proc.Tracks, track)
	return track
}

// ThreadByID gets the thread with the appropriate ID
func (proc *Proc) ThreadByID(tid trace.ID) (*Thread, bool) {
	for _, thread := range proc.Threads {
		if thread.TID == tid {
			return thread, true
		}
	}
	return nil, false
}

type Thread struct {
	TID trace.ID

	Start trace.Time
	Stop  trace.Time

	Layers []*Layer
}

func (thread *Thread) OpenLayer(time trace.Time) *Layer {
	for _, lay := range thread.Layers {
		if lay.End() <= time {
			return lay
		}
	}
	lay := &Layer{}
	thread.Layers = append(thread.Layers, lay)
	return lay
}

func (thread *Thread) CloseLayers(time trace.Time) {
	for _, lay := range thread.Layers {
		if len(lay.Spans) > 0 {
			span := &lay.Spans[len(lay.Spans)-1]
			if span.Stop.Unassigned() {
				span.Stop = time
			}
		}
	}
}

func (thread *Thread) LayerWithID(id trace.ID) (*Layer, bool) {
	for _, lay := range thread.Layers {
		if lay.OpenSpanID() == id {
			return lay, true
		}
	}
	return nil, false
}

type Track struct {
	Threads []*Thread
}

func (track *Track) Begin() trace.Time {
	if len(track.Threads) == 0 {
		return trace.MaxTime
	}
	return track.Threads[0].Start
}

func (track *Track) End() trace.Time {
	if len(track.Threads) == 0 {
		return trace.MinTime
	}
	return track.Threads[len(track.Threads)-1].Stop
}

type Span struct {
	ID    trace.ID
	Start trace.Time
	Stop  trace.Time
}

type Layer struct {
	Spans []Span
}

func (lay *Layer) lastSpan() *Span {
	if len(lay.Spans) == 0 {
		return nil
	}
	return &lay.Spans[len(lay.Spans)-1]
}

func (lay *Layer) Begin() trace.Time {
	if len(lay.Spans) == 0 {
		return trace.MaxTime
	}
	return lay.Spans[0].Start
}

func (lay *Layer) End() trace.Time {
	if span := lay.lastSpan(); span != nil {
		return span.Stop
	}
	return trace.MinTime
}

func (lay *Layer) OpenSpanID() trace.ID {
	if span := lay.lastSpan(); span != nil {
		if span.Stop.Unassigned() {
			return span.ID
		}
	}
	return trace.InvalidID
}

func (lay *Layer) OpenSpan(id trace.ID, time trace.Time) {
	assert(lay.End() <= time)

	lay.Spans = append(lay.Spans, Span{
		ID:    id,
		Start: time,
		Stop:  trace.MaxTime,
	})
}

func (lay *Layer) CloseSpan(id trace.ID, time trace.Time) {
	last := &lay.Spans[len(lay.Spans)-1]
	assert(last.ID == id)
	last.Stop = time
}
