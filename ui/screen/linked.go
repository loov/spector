package screen

import (
	"sort"

	"github.com/egonelbre/spector/ui/g"
)

func FindLinked(screen *Screen, at, relat g.Vector) (x, y *Segment) {
	xsegs := []*Segment{}
	ysegs := []*Segment{}

	// collect potential segments
	for _, area := range screen.Areas {
		// hit := area.Bounds.Test(at, AreaBorderRadius)

		if g.Abs(area.Bounds.Min.X-at.X) < AreaBorderRadius {
			ysegs = append(ysegs, &Segment{
				Min: area.RelBounds.Min.Y,
				Max: area.RelBounds.Max.Y,
				Ps:  []*float32{&area.RelBounds.Min.X},
			})
		}

		if g.Abs(area.Bounds.Min.Y-at.Y) < AreaBorderRadius {
			xsegs = append(xsegs, &Segment{
				Min: area.RelBounds.Min.X,
				Max: area.RelBounds.Max.X,
				Ps:  []*float32{&area.RelBounds.Min.Y},
			})
		}

		if g.Abs(area.Bounds.Max.X-at.X) < AreaBorderRadius {
			ysegs = append(ysegs, &Segment{
				Min: area.RelBounds.Min.Y,
				Max: area.RelBounds.Max.Y,
				Ps:  []*float32{&area.RelBounds.Max.X},
			})
		}

		if g.Abs(area.Bounds.Max.Y-at.Y) < AreaBorderRadius {
			xsegs = append(xsegs, &Segment{
				Min: area.RelBounds.Min.X,
				Max: area.RelBounds.Max.X,
				Ps:  []*float32{&area.RelBounds.Max.Y},
			})
		}
	}

	xsegs, ysegs = Merge(xsegs), Merge(ysegs)
	return FindMatch(xsegs, relat.X), FindMatch(ysegs, relat.Y)
}

func Merge(segs []*Segment) []*Segment {
	if len(segs) == 0 {
		return nil
	}
	sort.Slice(segs, func(i, k int) bool {
		return segs[i].Less(segs[k])
	})

	xs := []*Segment{segs[0]}
	for _, seg := range segs[1:] {
		if !xs[len(xs)-1].Merge(seg, AreaBorderRadius) {
			xs = append(xs, seg)
		}
	}

	return xs
}

func FindMatch(segs []*Segment, p float32) *Segment {
	for _, seg := range segs {
		if seg.Contains(p, AreaBorderRadius) {
			return seg
		}
	}
	return &Segment{Min: p, Max: p}
}

type Segment struct {
	Min, Max float32
	Ps       []*float32
}

func (a *Segment) Contains(p, r float32) bool {
	return a.Min-r <= p && p <= a.Max+r
}

func (a *Segment) Merge(b *Segment, r float32) bool {
	if a.Contains(b.Min, r) {
		if a.Max < b.Max {
			a.Max = b.Max
		}
		a.Ps = append(a.Ps, b.Ps...)
		return true
	}
	return false
}

func (a *Segment) Less(b *Segment) bool {
	if a.Min == b.Min {
		return a.Max > b.Max
	}
	return a.Min < b.Min
}
