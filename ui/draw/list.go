package draw

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
	list.CurrentClip = Rectangle{InvalidVector, InvalidVector}
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
