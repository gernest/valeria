package wuxia

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGenerator_Build(t *testing.T) {
	g := NewGenerator(nil, nil, afero.NewOsFs())
	err := g.SetWorkDir("fixture/site")
	if err != nil {
		t.Error(err)
	}
	err = g.Build()
	if err != nil {
		t.Error(err)
	}
}