package wuxia

import (
	"encoding/json"
	"errors"

	"github.com/robertkrimen/otto"
)

func getToJSON(o interface{}, vm *otto.Otto, src string) error {
	source := "JSON.stringify( function(){" + src + "}())"
	v, err := vm.Run(source)
	if err != nil {
		return err
	}
	if v.IsUndefined() || v.IsNull() {
		return errors.New("bad script")
	}
	s, err := v.ToString()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(s), o)
	if err != nil {
		return err
	}
	return nil
}
func Panic(o interface{}) {
	v, err := otto.ToValue(o)
	if err != nil {
		errV, _ := otto.ToValue(err)
		panic(errV)
	}
	panic(v)
}

func ToValue(o interface{}) otto.Value {
	v, _ := otto.ToValue(o)
	return v
}
