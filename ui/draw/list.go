package draw

var zeroClip = Rectangle{Vector{-8192, -8192}, Vector{+8192, +8192}}

type List struct {
	Commands []Command
	Indicies []Index
	Vertices []Vertex

	CurrentCommand *Command
	CurrentClip    Rectangle
	CurrentTexture TextureID

	ClipStack    []Rectangle
	TextureStack []TextureID
}

func NewList() *List {
	list := &List{}
	list.Reset()
	return list
}

func (list *List) Reset() {
	list.Commands = list.Commands[:0:cap(list.Commands)]
	list.Indicies = list.Indicies[:0:cap(list.Indicies)]
	list.Vertices = list.Vertices[:0:cap(list.Vertices)]

	list.CurrentCommand = nil
	list.CurrentClip = zeroClip
	list.CurrentTexture = 0

	list.ClipStack = nil
	list.TextureStack = nil

	list.BeginCommand()
}

func (list *List) PushClip(clip Rectangle) {
	list.ClipStack = append(list.ClipStack, list.CurrentClip)
	list.CurrentClip = clip
	list.updateClip()
}

func (list *List) PushClipFullscreen() { list.PushClip(zeroClip) }

func (list *List) PopClip() {
	n := len(list.ClipStack)
	list.CurrentClip = list.ClipStack[n-1]
	list.ClipStack = list.ClipStack[:n-1]
	list.updateClip()
}

func (list *List) updateClip() {
	if list.CurrentCommand == nil ||
		list.CurrentCommand.Clip != list.CurrentClip {
		list.BeginCommand()
		return
	}
	list.CurrentCommand.Clip = list.CurrentClip
}

func (list *List) PushTexture(id TextureID) {
	list.TextureStack = append(list.TextureStack, list.CurrentTexture)
	list.CurrentTexture = id
	list.updateTexture()
}

func (list *List) PopTexture() {
	n := len(list.TextureStack)
	list.CurrentTexture = list.TextureStack[n-1]
	list.TextureStack = list.TextureStack[:n-1]
	list.updateTexture()
}

func (list *List) updateTexture() {
	if list.CurrentCommand == nil ||
		list.CurrentCommand.Texture != list.CurrentTexture {
		list.BeginCommand()
		return
	}
	list.CurrentCommand.Texture = list.CurrentTexture
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
