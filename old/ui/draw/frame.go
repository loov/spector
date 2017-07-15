package draw

type Frame struct {
	Size  Vector
	Lists []*List
	lists []*List
}

func (frame *Frame) Reset() {
	frame.lists = frame.Lists
	frame.Lists = nil
}

func (frame *Frame) Next() bool {
	frame.Reset()
	return true
}

func (frame *Frame) Layer() *List {
	var list *List

	if len(frame.lists) > 0 {
		list = frame.lists[0]
		frame.lists = frame.lists[1:]
	} else {
		list = NewList()
	}
	list.Reset()
	frame.Lists = append(frame.Lists)
	return list
}
