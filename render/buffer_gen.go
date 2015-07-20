// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type op struct {
	name string
	code byte
	spec *ast.TypeSpec
}

type bycode []*op

func (a bycode) Len() int           { return len(a) }
func (a bycode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bycode) Less(i, j int) bool { return a[i].code < a[j].code }

func parseOps(f *ast.File) []*op {
	r := []*op{}
	for _, obj := range f.Scope.Objects {
		spec, ok := obj.Decl.(*ast.TypeSpec)
		if !ok || spec.Doc == nil || len(spec.Doc.List) == 0 {
			continue
		}

		code := 0
		for _, line := range spec.Doc.List {
			t := strings.TrimSpace(line.Text)
			if !strings.HasPrefix(t, "// op ") {
				continue
			}
			t = strings.TrimPrefix(t, "// op ")

			_, err := fmt.Sscanf(t, "0x%x", &code)
			if err == nil {
				break
			}
		}

		if code == 0 {
			panic("did not find code for " + spec.Name.Name)
		}
		r = append(r, &op{
			name: spec.Name.Name,
			code: byte(code),
			spec: spec,
		})
	}

	sort.Sort(bycode(r))
	return r
}

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "ops.go", nil, parser.ParseComments)
	if err != nil {
		panic(err.Error())
		return
	}

	ops := parseOps(f)

	out := &bytes.Buffer{}
	fmt.Fprintln(out, "// GENERATED CODE")
	fmt.Fprintln(out, "// DO NOT MODIFY")
	fmt.Fprintln(out, "package render")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, `import "unsafe"`)
	fmt.Fprintln(out, "")

	fmt.Fprintln(out, "type Op byte")
	fmt.Fprintln(out, "const (")
	fmt.Fprintln(out, "OpInvalid = Op(0x00)")
	for _, op := range ops {
		fmt.Fprintf(out, "Op%s = Op(0x%02x)\n", op.name, op.code)
	}
	fmt.Fprintln(out, ")")

	fmt.Fprintln(out, "func (op Op) Size() int {")
	fmt.Fprintln(out, "switch op {")
	for _, op := range ops {
		fmt.Fprintf(out, "case Op%[1]s : return int(unsafe.Sizeof(%[1]s{}))\n", op.name)
	}
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out, `panic("invalid op")`)
	fmt.Fprintln(out, "}")

	for _, op := range ops {
		fmt.Fprintf(out, `
			func (w *Buffer) %[1]s() *%[1]s { return (*%[1]s)(w.alloc(Op%[1]s, int(unsafe.Sizeof(%[1]s{})))) }
			func (r *Reader) %[1]s() *%[1]s { return (*%[1]s)(r.ptr()) }
		`, op.name)
	}

	data, err := format.Source(out.Bytes())
	if err != nil {
		fmt.Println(out.String())
		panic(err.Error())
		return
	}

	err = ioutil.WriteFile("buffer_ops.go", data, 0777)
	if err != nil {
		panic(err.Error())
		return
	}
}
