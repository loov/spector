// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
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
	Type  Type
}

type Type string

const (
	Int    = Type("Int")
	Byte   = Type("Byte")
	UTF8   = Type("UTF8")
	Blob   = Type("Blob")
	Values = Type("Values")
)

var ZeroValue = map[Type]string{
	Int:    "0",
	Byte:   "0",
	UTF8:   "''",
	Blob:   "new Uint8Array()",
	Values: "new Array()",
}

var Events = []*Event{
	{Name: "Invalid", Spec: "Code:Byte=0x00"},

	// start & stop stream
	{Name: "StreamStart", Spec: "Code:Byte=0x01 ProcID MachineID Timestamp CPUFrequency"},
	{Name: "StreamStop", Spec: "Code:Byte=0x02 TimeDiff"},

	// start & stop execution thread
	{Name: "ThreadStart", Spec: "Code:Byte=0x03 TimeDiff ThreadID StackID"},
	{Name: "ThreadStop", Spec: "Code:Byte=0x04 TimeDiff ThreadID StackID"},

	// begin & end a span
	{Name: "Begin", Spec: "Code:Byte=0x05 TimeDiff ThreadID StackID ID"},
	{Name: "End", Spec: "Code:Byte=0x06 TimeDiff ThreadID StackID ID"},

	// start & finish an arrow
	{Name: "Start", Spec: "Code:Byte=0x07 TimeDiff ThreadID StackID ID"},
	{Name: "Finish", Spec: "Code:Byte=0x08 TimeDiff ThreadID StackID ID"},

	// sample integer values
	{Name: "Sample", Spec: "Code:Byte=0x09 TimeDiff ThreadID StackID Values:Values"},
	// create a snapshot from an item
	{Name: "Snapshot", Spec: "Code:Byte=0x0A TimeDiff ThreadID StackID ID ContentKind:Byte Content:Blob"},

	// provide information about a specific ID
	{Name: "Info", Spec: "Code:Byte=0x0B ID Name:UTF8 ContentKind:Byte Content:Blob"},
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
		f.Type = Type(spec[s:e])
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
			if f.Type == "" {
				f.Type = Int
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

var Funcs = template.FuncMap{
	"zero": func(f Field) string {
		if f.Value != "" {
			return f.Value
		}
		return ZeroValue[f.Type]
	},
}

var JS = template.Must(template.New("").Funcs(Funcs).Parse(`
<script>
package("spector", function(){
	var Event = {};
	var EventByCode = {};

	{{ range $event := .Events }}
	Event.{{$event.Name}} = {{$event.Name}}Event;
	{{$event.Name}}Event.Code = {{$event.Code}};
	{{$event.Name}}Event.prototype.Code = {{$event.Code}};
	EventByCode[{{$event.Code}}] = {{$event.Name}}Event;
	function {{$event.Name}}Event(props){
		props = props !== undefined ? props : {};
		{{ range $field := $event.Fields }}
		this.{{$field.Name}} = props.{{$field.Name}} || {{ zero $field }};{{ end }}
	};

	{{$event.Name}}Event.prototype = {
		read: function(stream){ {{ range $field := $event.Fields }}
			this.{{$field.Name}} = stream.read{{$field.Type}}();{{ end }}
		},
		write: function(stream){ {{ range $field := $event.Fields }}{{if $field.Value}}
			stream.write{{$field.Type}}({{$field.Value}});{{else}}
			stream.write{{$field.Type}}(this.{{$field.Name}});{{ end }}{{ end }}
		}
	};
	{{ end }}

	var ContentKind = { {{ range $kind := .ContentKinds }}
		{{$kind.Name}}: {{$kind.Code}},{{ end }}
	};

	return {
		Version: {{ .Version }},
		Event: Event,
		EventByCode: EventByCode,
		ContentKind: ContentKind
	};
})
</script>
`))

func main() {
	var buf bytes.Buffer
	check(JS.Execute(&buf, map[string]interface{}{
		"Version":      Version,
		"Events":       Events,
		"ContentKinds": ContentKinds,
	}))

	check(ioutil.WriteFile(filepath.Join("spector", "protocol.html"), buf.Bytes(), 0777))
}
