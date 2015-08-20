// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const Version = 1

type Event struct {
	Name   string
	Spec   string
	Code   string
	Fields []Field
}

type Field struct {
	Name  string
	Value string
	Kind  Kind
}

type Kind string

const (
	Int    = Kind("Int")
	Byte   = Kind("Byte")
	UTF8   = Kind("UTF8")
	Blob   = Kind("Blob")
	Values = Kind("Values")
)

func (k Kind) Type() string {
	switch k {
	case Int:
		return "int32"
	case Byte:
		return "byte"
	case UTF8:
		return "string"
	case Blob:
		return "[]byte"
	case Values:
		return "[]Value"
	}
	return ""
}

var Events = []*Event{
	{Name: "Invalid", Spec: "Code:Byte=0x00"},

	// start & stop stream
	{Name: "StreamStart", Spec: "Code:Byte=0x01 ProcessID MachineID Time CPUFrequency"},
	{Name: "StreamStop", Spec: "Code:Byte=0x02 Time"},

	// start & stop execution thread
	{Name: "ThreadStart", Spec: "Code:Byte=0x03 Time ThreadID StackID"},
	{Name: "ThreadSleep", Spec: "Code:Byte=0x04 Time ThreadID StackID"},
	{Name: "ThreadWake", Spec: "Code:Byte=0x05 Time ThreadID StackID"},
	{Name: "ThreadStop", Spec: "Code:Byte=0x06 Time ThreadID StackID"},

	// begin & end a span
	{Name: "Begin", Spec: "Code:Byte=0x07 Time ThreadID StackID ID"},
	{Name: "End", Spec: "Code:Byte=0x08 Time ThreadID StackID ID"},

	// start & finish an arrow
	{Name: "Start", Spec: "Code:Byte=0x09 Time ThreadID StackID ID"},
	{Name: "Finish", Spec: "Code:Byte=0x0A Time ThreadID StackID ID"},

	// sample integer values
	// {Name: "Sample", Spec: "Code:Byte=0x0B Time ThreadID StackID Values:Values"},
	// create a snapshot from an item
	{Name: "Snapshot", Spec: "Code:Byte=0x0C Time ThreadID StackID ID ContentKind:Byte Content:Blob"},

	// provide information about a specific ID
	{Name: "Info", Spec: "Code:Byte=0x0D ID Name:UTF8 ContentKind:Byte Content:Blob"},
}

type ContentKind struct {
	Name string
	Code string
}

var ContentKinds = []*ContentKind{
	{Name: "Invalid", Code: "0x00"},

	{Name: "Thread", Code: "0x01"},
	{Name: "Stack", Code: "0x02"},

	// generic types
	{Name: "Text", Code: "0x10"},
	{Name: "JSON", Code: "0x11"},
	{Name: "BLOB", Code: "0x12"},
	{Name: "Image", Code: "0x13"},

	// user types
	{Name: "User", Code: "0x20"},
}

var Code = template.Must(template.New("").Parse(`
package trace

type Event interface {
	Code() byte
	ReadFrom(r *Reader)
	WriteTo(w *Writer)
}

func NewEventByCode(code byte) Event {
	switch code {
	{{ range $event := .Events }}
	case {{$event.Code}}:
		return &{{$event.Name}}{}
	{{ end }}
	}
	panic("unknown code")
}

{{ range $event := .Events }}
// code: {{$event.Code}}
type {{$event.Name}} struct {
	{{ range $field := $event.Fields }}
	{{$field.Name}} {{$field.Kind.Type}}
	{{ end }}
}
{{ end }}

{{ range $event := .Events }}
func (ev *{{$event.Name}}) Code() byte { return {{$event.Code}} }

func (ev *{{$event.Name}}) ReadFrom(r *Reader) {
	{{ range $field := $event.Fields }}
	ev.{{$field.Name}} = r.enc.Read{{$field.Kind}}();
	{{ end }}
}

func (ev *{{$event.Name}}) WriteTo(w *Writer) {
	{{ range $field := $event.Fields }}
	w.enc.Write{{$field.Kind}}(ev.{{$field.Name}});
	{{ end }}
}
{{ end }}
`))

func main() {
	var buf bytes.Buffer
	check(Code.Execute(&buf, map[string]interface{}{
		"Version":      Version,
		"Events":       Events,
		"ContentKinds": ContentKinds,
	}))

	bytes := buf.Bytes()
	rx := regexp.MustCompile(`(?m)[ \t\n]+$`)
	bytes = rx.ReplaceAll(bytes, []byte{})
	bytes, err := format.Source(bytes)
	check(err)

	check(ioutil.WriteFile(filepath.Join("trace", "events.go"), bytes, 0777))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseField(spec string) Field {
	f := Field{}

	s, e := 0, 0
	for e < len(spec) && spec[e] != ':' {
		e++
	}
	f.Name = spec[s:e]

	s, e = e+1, e+1
	for e < len(spec) && spec[e] != '=' {
		e++
	}

	if s < e {
		f.Kind = Kind(spec[s:e])
	}

	e = e + 1
	if e < len(spec) {
		f.Value = spec[e:]
	}

	return f
}

func init() {
	for _, ev := range Events {
		fields := strings.Fields(ev.Spec)
		for _, field := range fields {
			f := parseField(field)
			if f.Kind == "" {
				f.Kind = Int
			}
			ev.Fields = append(ev.Fields, f)
		}
		if ev.Fields[0].Name != "Code" {
			panic("invalid code")
		}
		ev.Code = ev.Fields[0].Value
		ev.Fields = ev.Fields[1:]
	}
}
