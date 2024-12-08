// Command jtog is a command line tool that converts JSON to Go source code.
//
// Usage: jtog [ -i=bool ] [ -l=bool ] [ -o=bool ] [ file ... ]
// If no file path(s) are specified as flags then data from standard
// input is assumed.
//
//	-i	indent using spaces
//	-l	inline type defintions (default true)
//	-o	appends "omitempty" to the json tag
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var (
	indentflag = flag.Bool("i", false, "indent using spaces")
	inlineflag = flag.Bool("l", true, "inline type defintions")
	omitflag   = flag.Bool("o", false, "appends \"omitempty\" to the json tag")

	sb strings.Builder
)

type Field struct {
	Name   string
	Type   string
	Tag    string
	Fields []Field
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [ -i=bool ] [ -l=bool ] [ -o=bool ] [ file ... ]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "If no file path(s) are specified as flags then data from standard input is assumed.\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, arg := range args {
		f := os.Stdin
		if arg != "-" {
			var err error
			f, err = os.Open(flag.Args()[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer f.Close()
		}

		buf, err := io.ReadAll(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(buf, &data); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		dump(parse(map[string]interface{}{filepath.Base(os.Args[0]): data}, *inlineflag, *omitflag), *inlineflag, 0)
		fmt.Print(sb.String())
	}
}

func emit(indent int, format string, a ...any) {
	ind := "	"
	if *indentflag {
		ind = "    "
	}
	for range indent {
		sb.WriteString(ind)
	}
	sb.WriteString(fmt.Sprintf(format, a...))
}

func dump(fields []Field, inline bool, indent int) {
	for _, f := range fields {
		if inline || indent == 0 {
			if indent == 0 {
				emit(indent, "type %s struct {\n", f.Name)
			} else {
				emit(indent, "%s struct {\n", f.Name)
			}
			indent++
			for _, sf := range f.Fields {
				if sf.Type == "struct" {
					dump([]Field{sf}, inline, indent)
					continue
				}
				emit(indent, "%s	%s	`json:\"%s\"`\n", sf.Name, sf.Type, sf.Tag)
			}
			indent--
			if indent == 0 {
				emit(indent, "}\n")
			} else {
				emit(indent, "} `json:\"%s\"`\n", f.Tag)
			}
			continue
		}
		emit(indent, "%s	%s	`json:\"%s\"`\n", f.Name, f.Type, f.Tag)
		buf, i := sb.String(), indent
		sb.Reset()
		indent = 0
		dump([]Field{f}, inline, indent)
		newstruct := sb.String()
		sb.Reset()
		sb.WriteString(newstruct + "\n")
		sb.WriteString(buf)
		indent = i
	}
}

func parse(data map[string]interface{}, inline, omitempty bool) []Field {
	var fields []Field
	tag := ",omitempty"
	if !omitempty {
		tag = ""
	}
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			f := parse(v, inline, omitempty)
			fields = append(fields, Field{Name: strings.Title(key), Type: "struct", Tag: key + tag, Fields: f})
		case []interface{}:
			typ := "any"
			if len(v) > 0 {
				typ = fmt.Sprint(reflect.TypeOf(v[0]))
			}
			fields = append(fields, Field{Name: strings.Title(key), Type: fmt.Sprintf("[]%s", typ), Tag: key + tag})
		case float64:
			typ := "float64"
			if _, err := strconv.Atoi(fmt.Sprint(value)); err == nil {
				typ = "int"
			} else if _, err := strconv.ParseInt(fmt.Sprint(value), 10, 64); err == nil {
				typ = "int64"
			}
			fields = append(fields, Field{Name: strings.Title(key), Type: typ, Tag: key + tag})
		case bool:
			fields = append(fields, Field{Name: strings.Title(key), Type: "bool", Tag: key + tag})
		case string:
			fields = append(fields, Field{Name: strings.Title(key), Type: "string", Tag: key + tag})
		case int:
			fields = append(fields, Field{Name: strings.Title(key), Type: "int", Tag: key + tag})
		}
	}
	return fields
}
