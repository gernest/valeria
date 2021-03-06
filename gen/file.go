package gen

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/spf13/afero"
)

//Export helper for building javascript objects from Go objects.
type Export map[string]interface{}

//Set add a new key key, to the export object with value value.
func (e Export) Set(key string, value interface{}) {
	e[key] = value
}

//ToValue converts the export object to otto.Value
func (e Export) ToValue(vm *otto.Otto) otto.Value {
	o, err := vm.Object(`({})`)
	if err != nil {
		panicOtto(err.Error())
	}
	for key, value := range e {
		_ = o.Set(key, value)
	}
	return o.Value()
}

const (
	fileFlagRead   = "r"
	fileFlagWrite  = "w"
	fileFlagCreate = "c"
	fileFlagAppend = "a"
	fileFlagTrucc  = "t"
)

type fileSys struct {
	afero.Fs
}

func (fs fileSys) export() Export {
	e := make(Export)
	e.Set("open", fs.open)
	e.Set("openFile", fs.openFile)
	return e
}

func (fs fileSys) open(call otto.FunctionCall) otto.Value {
	name, err := call.Argument(0).ToString()
	if err != nil {
		panicOtto(err.Error())
	}
	f, err := fs.Open(name)
	if err != nil {
		panicOtto(err.Error())
	}
	af := &file{o: f}
	return af.export().ToValue(call.Otto)
}

func buildFlags(src string) (int, error) {
	parts := strings.Split(src, "+")
	if len(parts) > 0 {
		var f int
		for i := 0; i < len(parts); i++ {
			switch parts[i] {
			case fileFlagRead:
				f = f | os.O_RDONLY
			case fileFlagWrite:
				f = f | os.O_WRONLY
			case fileFlagCreate:
				f = f | os.O_CREATE
			case fileFlagTrucc:
				f = f | os.O_TRUNC
			default:
				return f, errors.New("unknown flag " + parts[i])
			}
		}
		return f, nil
	}
	return 0, errors.New("no flags found")
}

func (fs fileSys) openFile(call otto.FunctionCall) otto.Value {
	name, err := call.Argument(0).ToString()
	if err != nil {
		panicOtto(err.Error())
	}
	flag, err := call.Argument(1).ToString()
	if err != nil {
		panicOtto(err.Error())
	}
	uflag, err := buildFlags(flag)
	if err != nil {
		panicOtto(err.Error())
	}
	mode, err := call.Argument(2).ToString()
	if err != nil {
		panicOtto(err.Error())
	}
	umode, err := strconv.ParseUint(mode, 10, 32)
	if err != nil {
		panicOtto(err.Error())
	}
	f, err := fs.OpenFile(name, uflag, os.FileMode(umode))
	if err != nil {
		panicOtto(err.Error())
	}
	af := &file{o: f}
	return af.export().ToValue(call.Otto)
}

type file struct {
	o afero.File
}

func (f *file) export() Export {
	e := make(Export)
	e.Set("close", f.close)
	e.Set("read", f.read)
	e.Set("write", f.write)
	return e
}

func (f *file) close(call otto.FunctionCall) otto.Value {
	err := f.o.Close()
	if err != nil {
		panicOtto(err.Error())
	}
	return otto.Value{}
}

func (f *file) read(call otto.FunctionCall) otto.Value {
	b, err := ioutil.ReadAll(f.o)
	if err != nil {
		panicOtto(err.Error())
	}
	return ToValue(string(b))
}

func (f *file) write(call otto.FunctionCall) otto.Value {
	c, err := call.Argument(0).ToString()
	if err != nil {
		panicOtto(err.Error())
	}
	n, err := f.o.WriteString(c)
	if err != nil {
		panicOtto(err.Error())
	}
	return ToValue(n)
}
