package timeline

import (
	"github.com/egonelbre/spector/trace"
)

type Handler struct {
	Timeline *Timeline
	Proc     *Proc
}

func (h *Handler) Handle(ev trace.Event) {
	timeline, proc := h.Timeline, h.Proc
	timeline.TotalEvents++

	switch ev := ev.(type) {
	case *trace.StreamStart:
		h.Proc = &Proc{
			PID:   ev.ProcessID,
			MID:   ev.MachineID,
			Time:  ev.Time,
			Freq:  ev.Freq,
			Start: ev.Time,
			Stop:  trace.MaxTime,
		}
		timeline.Procs = append(timeline.Procs, h.Proc)
	case *trace.StreamStop:
		assert(proc != nil)

		proc.Stop = ev.Time
		proc = nil

	case *trace.ThreadStart:
		assert(proc != nil)

		thread := &Thread{
			TID:   ev.ThreadID,
			Start: ev.Time,
			Stop:  trace.MaxTime,
		}

		track := proc.OpenTrack(ev.Time)
		track.Threads = append(track.Threads, thread)
		proc.Threads = append(proc.Threads, thread)
	case *trace.ThreadStop:
		assert(proc != nil)
		thread, ok := proc.ThreadByID(ev.ThreadID)
		assert(ok)

		thread.Stop = ev.Time
		thread.CloseLayers(ev.Time)
	case *trace.ThreadSleep:
		assert(proc != nil)
		_, ok := proc.ThreadByID(ev.ThreadID)
		assert(ok)

		//TODO
	case *trace.ThreadWake:
		assert(proc != nil)
		_, ok := proc.ThreadByID(ev.ThreadID)
		assert(ok)
		//TODO

	case *trace.SpanBegin:
		assert(proc != nil)
		thread, ok := proc.ThreadByID(ev.ThreadID)
		assert(ok)

		layer := thread.OpenLayer(ev.Time)
		layer.OpenSpan(ev.ID, ev.Time)
	case *trace.SpanEnd:
		assert(proc != nil)
		thread, ok := proc.ThreadByID(ev.ThreadID)
		assert(ok)

		layer, ok := thread.LayerWithID(ev.ID)
		assert(ok)
		layer.CloseSpan(ev.ID, ev.Time)
	}
}
