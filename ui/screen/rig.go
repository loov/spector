package screen

import (
	"sort"

	"github.com/egonelbre/spector/ui/g"
)

var (
	RigJoinerColor          = g.Color{0xA0, 0xA0, 0xA0, 0xFF}
	RigJoinerHighlightColor = g.Color{0xE0, 0xA0, 0xA0, 0xFF}
	RigBorderColor          = g.Color{0x80, 0x80, 0x80, 0xFF}
	RigBorderHighlightColor = g.Color{0xA0, 0x80, 0x80, 0xFF}
	RigCornerColor          = g.Color{0x40, 0x40, 0x40, 0xFF}
	RigCornerHighlightColor = g.Color{0xA0, 0x40, 0x40, 0xFF}

	RigBackground = g.Color{0x80, 0x80, 0x80, 0xFF}

	RigJoinerSize    = float32(20)
	RigTriggerSize   = float32(40)
	RigTriggerRadius = g.Vector{RigTriggerSize, RigTriggerSize}
	RigBorderRadius  = g.Vector{2, 2}
	RigCornerRadius  = g.Vector{3, 3}
)

type Rig struct {
	Corners []*Corner
	Borders []*Border
}

func NewRig() *Rig {
	topLeft := &Corner{Pos: g.Vector{0, 0}}
	topRight := &Corner{Pos: g.Vector{1, 0}}
	bottomLeft := &Corner{Pos: g.Vector{0, 1}}
	bottomRight := &Corner{Pos: g.Vector{1, 1}}

	leftSide := &Border{Locked: true, Horizontal: false, Corners: []*Corner{bottomLeft, topLeft}}
	topSide := &Border{Locked: true, Horizontal: true, Corners: []*Corner{topLeft, topRight}}
	rightSide := &Border{Locked: true, Horizontal: false, Corners: []*Corner{topRight, bottomRight}}
	bottomSide := &Border{Locked: true, Horizontal: true, Corners: []*Corner{bottomRight, bottomLeft}}

	leftSide.Sort()
	topSide.Sort()
	rightSide.Sort()
	bottomSide.Sort()

	topLeft.Vertical = leftSide
	topLeft.Horizontal = topSide
	topRight.Horizontal = topSide
	topRight.Vertical = rightSide
	bottomRight.Vertical = rightSide
	bottomRight.Horizontal = bottomSide
	bottomLeft.Horizontal = bottomSide
	bottomLeft.Vertical = leftSide

	return &Rig{
		Corners: []*Corner{topLeft, topRight, bottomLeft, bottomRight},
		Borders: []*Border{leftSide, topSide, rightSide, bottomSide},
	}
}

func (rig *Rig) SplitVertically(topRight *Corner, posX float32) *Border {
	_, bottomRight := topRight.BlockingVertical(true, false)
	topSide := topRight.Horizontal
	bottomSide := bottomRight.Horizontal

	centerTop := &Corner{Pos: g.Vector{X: posX, Y: topRight.Pos.Y}}
	centerBottom := &Corner{Pos: g.Vector{X: posX, Y: bottomRight.Pos.Y}}
	split := &Border{Locked: false, Horizontal: false, Corners: []*Corner{centerTop, centerBottom}}
	split.Sort()

	centerTop.Horizontal, centerTop.Vertical = topSide, split
	topSide.Insert(centerTop)

	centerBottom.Horizontal, centerBottom.Vertical = bottomSide, split
	bottomSide.Insert(centerBottom)

	rig.Corners = append(rig.Corners, split.Corners...)
	SortCorners(rig.Corners)
	rig.Borders = append(rig.Borders, split)

	return split
}

func (rig *Rig) SplitHorizontally(topRight *Corner, posY float32) *Border {
	topLeft, _ := topRight.BlockingHorizontal(false, true)
	rightSide := topRight.Vertical
	leftSide := topLeft.Vertical

	centerLeft := &Corner{Pos: g.Vector{X: topLeft.Pos.X, Y: posY}}
	centerRight := &Corner{Pos: g.Vector{X: topRight.Pos.X, Y: posY}}

	split := &Border{Locked: false, Horizontal: true, Corners: []*Corner{centerLeft, centerRight}}
	split.Sort()

	centerLeft.Vertical, centerLeft.Horizontal = leftSide, split
	leftSide.Insert(centerLeft)

	centerRight.Vertical, centerRight.Horizontal = rightSide, split
	rightSide.Insert(centerRight)

	rig.Corners = append(rig.Corners, split.Corners...)
	SortCorners(rig.Corners)
	rig.Borders = append(rig.Borders, split)

	return split
}

type Corner struct {
	Pos g.Vector

	Horizontal *Border
	Vertical   *Border
}

func (a *Corner) IsLocked() bool {
	return a.Horizontal.IsLocked() || a.Vertical.IsLocked()
}

func (a *Corner) SideLeft() *Border {
	if a == nil || a.Horizontal == nil || a.Horizontal.First() == a {
		return nil
	}
	return a.Horizontal
}
func (a *Corner) SideTop() *Border {
	if a == nil || a.Vertical == nil || a.Vertical.First() == a {
		return nil
	}
	return a.Vertical
}
func (a *Corner) SideRight() *Border {
	if a == nil || a.Horizontal == nil || a.Horizontal.Last() == a {
		return nil
	}
	return a.Horizontal
}
func (a *Corner) SideBottom() *Border {
	if a == nil || a.Vertical == nil || a.Vertical.Last() == a {
		return nil
	}
	return a.Vertical
}

func (a *Corner) CornerLeft() *Corner   { return a.SideLeft().Neighbor(a, -1) }
func (a *Corner) CornerTop() *Corner    { return a.SideTop().Neighbor(a, -1) }
func (a *Corner) CornerRight() *Corner  { return a.SideRight().Neighbor(a, 1) }
func (a *Corner) CornerBottom() *Corner { return a.SideBottom().Neighbor(a, 1) }

func (a *Corner) BlockingHorizontal(checkTop, checkBottom bool) (left, right *Corner) {
	index := a.Horizontal.Index(a)
	neighbors := a.Horizontal.Corners

	for k := index - 1; k >= 0; k-- {
		n := neighbors[k]
		if checkTop && n.SideTop() != nil {
			left = n
			break
		}
		if checkBottom && n.SideBottom() != nil {
			left = n
			break
		}
	}

	for k := index + 1; k < len(neighbors); k++ {
		n := neighbors[k]
		if checkTop && n.SideTop() != nil {
			right = n
			break
		}
		if checkBottom && n.SideBottom() != nil {
			right = n
			break
		}
	}
	return
}

func (a *Corner) BlockingVertical(checkLeft, checkRight bool) (top, bottom *Corner) {
	index := a.Vertical.Index(a)
	neighbors := a.Vertical.Corners

	for k := index - 1; k >= 0; k-- {
		n := neighbors[k]
		if checkLeft && n.SideLeft() != nil {
			top = n
			break
		}
		if checkRight && n.SideRight() != nil {
			top = n
			break
		}
	}

	for k := index + 1; k < len(neighbors); k++ {
		n := neighbors[k]
		if checkLeft && n.SideLeft() != nil {
			bottom = n
			break
		}
		if checkRight && n.SideRight() != nil {
			bottom = n
			break
		}
	}
	return
}

func (a *Corner) Less(b *Corner) bool {
	if a.Pos.X != b.Pos.X {
		return a.Pos.X < b.Pos.X
	}
	return a.Pos.Y < b.Pos.Y
}

type Border struct {
	Locked     bool
	Horizontal bool
	Corners    []*Corner
}

func (border *Border) IsLocked() bool { return (border != nil) && border.Locked }

func (border *Border) Index(corner *Corner) int {
	if border == nil {
		return -1
	}
	for i, c := range border.Corners {
		if c == corner {
			return i
		}
	}
	return -1
}

func (border *Border) Insert(corner *Corner) {
	border.Corners = append(border.Corners, corner)
	border.Sort()
}

func (border *Border) Neighbor(corner *Corner, di int) *Corner {
	if border == nil {
		return nil
	}

	i := border.Index(corner)
	if i < 0 {
		return nil
	}

	ti := i + di
	if 0 <= ti && ti < len(border.Corners) {
		return border.Corners[ti]
	}

	return nil
}

func (border *Border) Center() g.Vector { return border.Min().Add(border.Max()).Scale(0.5) }

func (border *Border) First() *Corner { return border.Corners[0] }
func (border *Border) Last() *Corner  { return border.Corners[len(border.Corners)-1] }

func (border *Border) Min() g.Vector { return border.First().Pos }
func (border *Border) Max() g.Vector { return border.Last().Pos }

func (border *Border) Sort() { SortCorners(border.Corners) }

func SortCorners(corners []*Corner) {
	sort.Slice(corners, func(i, k int) bool {
		return corners[i].Less(corners[k])
	})
}
