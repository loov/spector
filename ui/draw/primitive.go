package draw

func (list *List) BeginCommand() {
	list.Commands = append(list.Commands, Command{
		Clip:    list.CurrentClip,
		Texture: list.CurrentTexture,
	})
	list.CurrentCommand = &list.Commands[len(list.Commands)-1]
}

func (list *List) Primitive_Reserve(index_count, vertex_count int) {
	cmd := list.CurrentCommand
	cmd.Count += Index(index_count)
}

func (list *List) Primitive_Rect(r *Rectangle, color Color) {
	a, b, c, d := r.Corners()

	base := Index(len(list.Vertices))
	list.Indicies = append(list.Indicies,
		base+0, base+1, base+2,
		base+0, base+2, base+3,
	)
	list.Vertices = append(list.Vertices,
		Vertex{a, noUV, color},
		Vertex{b, noUV, color},
		Vertex{c, noUV, color},
		Vertex{d, noUV, color},
	)
}

func (list *List) Primitive_RectUV(r *Rectangle, uv *Rectangle, color Color) {
	a, b, c, d := r.Corners()
	uv_a, uv_b, uv_c, uv_d := uv.Corners()

	base := Index(len(list.Vertices))
	list.Indicies = append(list.Indicies,
		base+0, base+1, base+2,
		base+0, base+2, base+3,
	)
	list.Vertices = append(list.Vertices,
		Vertex{a, uv_a, color},
		Vertex{b, uv_b, color},
		Vertex{c, uv_c, color},
		Vertex{d, uv_d, color},
	)
}

func (list *List) Primitive_Quad(a, b, c, d Vector, color Color) {
	base := Index(len(list.Vertices))
	list.Indicies = append(list.Indicies,
		base+0, base+1, base+2,
		base+0, base+2, base+3,
	)
	list.Vertices = append(list.Vertices,
		Vertex{a, noUV, color},
		Vertex{b, noUV, color},
		Vertex{c, noUV, color},
		Vertex{d, noUV, color},
	)
}

func (list *List) Primitive_QuadUV(q *[4]Vector, uv *Rectangle, color Color) {
	a, b, c, d := q[0], q[1], q[2], q[3]
	uv_a, uv_b, uv_c, uv_d := uv.Corners()

	base := Index(len(list.Vertices))
	list.Indicies = append(list.Indicies,
		base+0, base+1, base+2,
		base+0, base+2, base+3,
	)
	list.Vertices = append(list.Vertices,
		Vertex{a, uv_a, color},
		Vertex{b, uv_b, color},
		Vertex{c, uv_c, color},
		Vertex{d, uv_d, color},
	)
}

func linenormal(a, b Vector, size float32) Vector {
	d := b.Sub(a)
	n := Vector{-d.Y, d.X}
	return n.ScaleTo(size)
}

func (list *List) AddLine(points []Vector, closed bool, thickness float32, color Color) {
	if len(points) < 2 || color.Transparent() {
		return
	}
	if closed && len(points) < 3 {
		closed = false
	}

	if !closed {
		qthick := thickness * 0.25

		list.Primitive_Reserve(
			(len(points)-1)*6,
			(len(points)-1)*4,
		)

		a, b := points[0], points[1]
		abn := linenormal(a, b, qthick)
		a1, a2 := a.Add(abn.Scale(2.0)), a.Sub(abn.Scale(2.0))
		for _, c := range points[2:] {
			bcn := linenormal(b, c, qthick)

			z := abn.Add(bcn)
			b1, b2 := b.Add(z), b.Sub(z)

			list.Primitive_Quad(a1, b1, b2, a2, color)

			abn = bcn
			a, b = b, c
			a1, a2 = b1, b2
		}

		b1, b2 := b.Add(abn.Scale(2.0)), b.Sub(abn.Scale(2.0))
		list.Primitive_Quad(a1, b1, b2, a2, color)
	} else {
		qthick := thickness * 0.25

		const X = 0
		list.Primitive_Reserve(
			(len(points)-X)*6,
			(len(points)-X)*4,
		)

		w := points[len(points)-3]
		a, b := points[len(points)-2], points[len(points)-1]
		wan := linenormal(w, a, qthick)
		abn := linenormal(a, b, qthick)
		z := wan.Add(abn)
		a1, a2 := a.Add(z), a.Sub(z)
		for _, c := range points[:len(points)-X] {
			bcn := linenormal(b, c, qthick)

			z := abn.Add(bcn)
			b1, b2 := b.Add(z), b.Sub(z)

			list.Primitive_Quad(a1, b1, b2, a2, color)

			abn = bcn
			a, b = b, c
			a1, a2 = b1, b2
		}

		/*		c := points[0]
				bcn := linenormal(b, c, qthick)
				z = abn.Add(bcn)
				b1, b2 := b.Add(z), b.Sub(z)

				list.Primitive_Quad(a1, b1, b2, a2, color)*/
	}
}

func (list *List) AddRectFill(r *Rectangle, color Color) {
	if color.Transparent() {
		return
	}

	list.Primitive_Reserve(6, 4)
	list.Primitive_Rect(r, color)
}
