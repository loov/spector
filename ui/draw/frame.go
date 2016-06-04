package draw

type Frame struct {
	Valid bool

	Text     *List
	Hint     *List
	Geometry *List
	Shadow   *List
}
