package main

import (
	"fmt"
	"redis"
	"reflect"
	"path/filepath"
	"strings"
//	"os"
	"io"
	"bytes"
	"text/template"
	"io/ioutil"
	"regexp"
)

var filter1, filter2, filter3 *regexp.Regexp
var LF []byte = []byte("\n")
func init() {
	var e error
	filter1, e = regexp.Compile(" *\n"); if e != nil {
		panic(e)
	}
	filter2, e = regexp.Compile("\t*\n"); if e != nil {
		panic(e)
	}
	filter3, e = regexp.Compile("\n{2,}\n"); if e != nil {
		panic(e)
	}
}
func main() {
	fmt.Printf("Salaam!\n")

	for _, ctype := range getSubjects() {
		fmt.Printf("\tctype: %s\n", ctype)

		t := loadMethodTestTemplate(ctype)
//		fmt.Printf("\ttemplate: %s\n", t)

		for i := 0; i < ctype.NumMethod(); i++ {
			m := ctype.Method(i)
//			renderFeatureTest(os.Stdout, ctype, m, t)
			var buf bytes.Buffer
			renderFeatureTest(&buf, ctype, m, t)
			writeTestFile(buf.Bytes(), ctype, m)
		}
	}
}

// REVU - conf this
var testdir = ".."

func loadMethodTestTemplate (subject reflect.Type) *template.Template {
//	dir := subject.Name()
	tfn := fmt.Sprintf("%s_mtest.tmpl", subject.Name())
	fn := filepath.Join(tfn)
	t, e := template.New(tfn).Funcs(funcmap()).ParseFiles(fn)
	if e != nil {
		fmt.Printf("Error creating template for %s - %s",subject, e)
		panic(e)
	}
	return t
}

func writeTestFile (data []byte, subject reflect.Type, method reflect.Method) {
	typedir := subject.Name()
	data = filter1.ReplaceAll(data, LF)
	data = filter2.ReplaceAll(data, LF)
	data = filter3.ReplaceAll(data, LF)
	fname := filepath.Join(testdir, typedir, fmt.Sprintf("%s_test.go", method.Name))
	if e := ioutil.WriteFile(fname, data, 0644); e != nil {
		panic(fmt.Errorf("error on WriteFile (%s) - %s", fname, e))
	}
}

type mtest struct {
	Type reflect.Type
	Subject string
	Method string
	InArgs   []string
	InCnt int
	OutArgs  []string
	OutCnt int
	Spec *redis.MethodSpec
}

func getMethodInArgs(m reflect.Method) []string {
	mt := m.Type
	args := make([]string, mt.NumIn())
	for i := 0; i<mt.NumIn(); i++ {
		args[i] = fmt.Sprintf("%s", mt.In(i))
	}
	return args
}

func getMethodOutArgs(m reflect.Method) []string {
	mt := m.Type
	args := make([]string, mt.NumOut())
	for i := 0; i<mt.NumOut(); i++ {
		args[i] = fmt.Sprintf("%s", mt.Out(i))
	}
	return args
}

func renderFeatureTest(w io.Writer, ctype reflect.Type, m reflect.Method, t *template.Template) {
	ins := getMethodInArgs(m)
	outs := getMethodOutArgs(m)
	data := &mtest{
		Type: ctype,
		Subject: ctype.Name(),
		Method:  m.Name,
		InArgs: ins,
		InCnt:  len(ins),
		OutArgs: outs,
		OutCnt: len(outs),
		Spec: redis.GetMethodSpec(ctype.Name(), m.Name),
	}
	t.Execute(w, data)
}

func getSubjects() []reflect.Type {
	return []reflect.Type {
		reflect.TypeOf((*redis.Client)(nil)).Elem(),
		reflect.TypeOf((*redis.AsyncClient)(nil)).Elem(),
	}
}

func funcmap() map[string]interface{} {
	return template.FuncMap {
		"zvTest": zvTest,
		"comma": comma,
		"isRedisError": isRedisError,
		"isFuture": isFuture,
		"isQuit": isQuit,
	}
}

func zvTest(tname string) string {
	switch tname {
	case "int64": return " != 0"
	case "float64": return " != float64(0)"
	case "bool": return " != false"
	case "string": return " != \"\""

	case "redis.KeyType": return " != redis.KeyType(0)"
	}
	return " != nil"
}

func comma(idx, cnt int) string {
	if idx < cnt -1 {
		return ", "
	}
	return ""
}

func isRedisError(v string) string {
	if v == "redis.Error" {
		return "true"
	}
	return ""
}

func isFuture(v string) string {
	if strings.HasPrefix(v, "redis.Future") {
		return "true"
	}
	return ""
}

func isQuit(v string) string {
	if v == "Quit" {
		return "true"
	}
	return ""
}
