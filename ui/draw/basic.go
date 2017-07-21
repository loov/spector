package draw

import "github.com/egonelbre/spector/ui/g"

func (list *List) StrokeLine(points []g.Vector, thickness float32, color g.Color) {
	if len(points) < 2 || color.Transparent() || thickness == 0 {
		return
	}

	startIndexCount := len(list.Indicies)

	R := g.Abs(thickness / 2.0)

	a := points[0]
	var x1, x2, xn g.Vector

	// draw each segment, where
	// a1-------^---------b1
	// |        | abn      |
	// a - - - - - - - - - b
	// |                   |
	// a2-----------------b2
	// x1, x2, xn are the previous segments end corners and normal
	for i, b := range points[1:] {
		// segment normal
		abn := g.SegmentNormal(a, b).ScaleTo(R)
		// segment corners
		a1, a2 := a.Add(abn), a.Sub(abn)
		b1, b2 := b.Add(abn), b.Sub(abn)

		if i > 0 && R > 1.5 {
			// draw segment chamfer
			d := xn.Rotate().Dot(abn)
			if d < 0 {
				list.Primitive_Tri(x1, a1, a, color)
			} else if d > 0 {
				list.Primitive_Tri(x2, a2, a, color)

			}
		}
		// draw block segment
		list.Primitive_Quad(a1, b1, b2, a2, color)

		a = b
		x1, x2, xn = b1, b2, abn
	}

	list.CurrentCommand.Count += Index(len(list.Indicies) - startIndexCount)
}

func (list *List) StrokeClosedLine(points []g.Vector, thickness float32, color g.Color) {
	if len(points) < 2 || color.Transparent() || thickness == 0 {
		return
	}
	if len(points) < 3 {
		list.StrokeLine(points, thickness, color)
		return
	}

	startIndexCount := len(list.Indicies)

	R := g.Abs(thickness / 2.0)
	a := points[len(points)-1]
	xn := g.SegmentNormal(points[len(points)-2], a).ScaleTo(R)
	x1, x2 := a.Add(xn), a.Sub(xn)

	// draw each segment, where
	// a1-------^---------b1
	// |        | abn      |
	// a - - - - - - - - - b
	// |                   |
	// a2-----------------b2
	// x1, x2, xn are the previous segments end corners and normal
	for _, b := range points {
		// segment normal
		abn := g.SegmentNormal(a, b).ScaleTo(R)
		// segment corners
		a1, a2 := a.Add(abn), a.Sub(abn)
		b1, b2 := b.Add(abn), b.Sub(abn)

		// draw segment chamfer
		if R > 1.5 {
			d := xn.Rotate().Dot(abn)
			if d < 0 {
				list.Primitive_Tri(x1, a1, a, color)
			} else if d > 0 {
				list.Primitive_Tri(x2, a2, a, color)
			}
		}

		// draw block segment
		list.Primitive_Quad(a1, b1, b2, a2, color)

		a = b
		x1, x2, xn = b1, b2, abn
	}

	list.CurrentCommand.Count += Index(len(list.Indicies) - startIndexCount)
}

const segmentsPerArc = 24

func (list *List) FillArc(center g.Vector, R float32, start, sweep float32, color g.Color) {
	if color.Transparent() || R == 0 {
		return
	}
	R = g.Abs(R)
	startIndexCount := len(list.Indicies)

	// N := sweep * R gives one segment per pixel
	N := Index(g.Clamp(g.Abs(sweep)*R/g.Tau, 3, segmentsPerArc))

	theta := sweep / float32(N)
	rots, rotc := g.Sincos(theta)
	dy, dx := g.Sincos(start)
	dy *= R
	dx *= R

	// add center point to the vertex buffer
	base := Index(len(list.Vertices))
	list.Vertices = append(list.Vertices, Vertex{center, NoUV, color})
	// add the first point the vertex buffer
	p := g.Vector{center.X + dx, center.Y + dy}
	list.Vertices = append(list.Vertices, Vertex{p, NoUV, color})
	// loop over rest of the points
	for i := Index(0); i < N; i++ {
		dx, dy = dx*rotc-dy*rots, dx*rots+dy*rotc
		p = g.Vector{center.X + dx, center.Y + dy}
		list.Vertices = append(list.Vertices, Vertex{p, NoUV, color})
		list.Indicies = append(list.Indicies, base, base+i+1, base+i+2)
	}

	list.CurrentCommand.Count += Index(len(list.Indicies) - startIndexCount)
}

func (list *List) FillCircle(center g.Vector, R float32, color g.Color) {
	if color.Transparent() || R == 0 {
		return
	}
	R = g.Abs(R)
	startIndexCount := len(list.Indicies)

	// N := 2 * PI * R gives one segment per pixel
	N := Index(g.Clamp(R, 3, segmentsPerArc))

	theta := g.Tau / float32(N)
	rots, rotc := g.Sincos(theta)

	dx, dy := R, float32(0)

	// add center point to the vertex buffer
	base := Index(len(list.Vertices))
	list.Vertices = append(list.Vertices, Vertex{center, NoUV, color})
	// add the first point the vertex buffer
	p := g.Vector{center.X + dx, center.Y + dy}
	list.Vertices = append(list.Vertices, Vertex{p, NoUV, color})

	// loop over rest of the points
	for i := Index(0); i < N; i++ {
		dx, dy = dx*rotc-dy*rots, dx*rots+dy*rotc
		p = g.Vector{center.X + dx, center.Y + dy}
		list.Vertices = append(list.Vertices, Vertex{p, NoUV, color})
		list.Indicies = append(list.Indicies, base, base+i+1, base+i+2)
	}

	list.CurrentCommand.Count += Index(len(list.Indicies) - startIndexCount)
}
