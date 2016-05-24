package draw

var zeroClip = Rectangle{Vector{-8192, -8192}, Vector{+8192, +8192}}

type List struct {
	Commands []Command
	Indicies []Index
	Vertices []Vertex

	CurrentCommand *Command
	CurrentClip    Rectangle
	CurrentTexture TextureID
}

func (list *List) Reset() {
	list.Commands = list.Commands[:0]
	list.Indicies = list.Indicies[:0]
	list.Vertices = list.Vertices[:0]

	list.CurrentCommand = nil
	list.CurrentClip = zeroClip
	list.CurrentTexture = 0
}

type Channel struct {
	Commands []Command
	Indicies []Index
}

type TextureID int32
type Callback func(*List, *Command)

type Command struct {
	Count    Index
	Clip     Rectangle
	Texture  TextureID
	Callback Callback
	Data     interface{}
}

type Index uint16

type Vertex struct {
	P     Vector
	UV    Vector
	Color Color
}
