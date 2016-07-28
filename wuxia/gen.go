package wuxia

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/robertkrimen/otto"
	"github.com/spf13/afero"
)

const (
	scriptsDir = "scripts"
	initDir    = "init"
	planDir    = "plan"
	configFile = "config.json"
	indexFile  = "index.js"
)

//BuildStage the step in the generation process life cycle.
type BuildStage int

// available stages offered by the generator
const (
	StageInit BuildStage = iota
	StageConfig
	StagePlan
	StageExec
)

func (s BuildStage) String() string {
	var rst string
	switch s {
	case StageInit:
		rst = "init"
	case StageConfig:
		rst = "config"
	case StagePlan:
		rst = "plan"
	case StageExec:
		rst = "exec"
	default:
		rst = "unkown stage"
	}
	return rst
}

//buildError error returned when building the static website. The error string
//returned is a json string that encodes the build stage and the message.
type buildError struct {
	Stage   string `json:"stage"`
	Message string `json:"msg"`
}

func buildErr(stage BuildStage, msg string) error {
	return &buildError{Stage: stage.String(), Message: msg}
}

func (b *buildError) Error() string {
	o, err := json.Marshal(b)
	if err != nil {
		return err.Error()
	}
	return string(o)
}

//Generator provides the static website generation capabilities.This is heavily
//integrated with the otto javascript runtime.
type Generator struct {
	vm  *otto.Otto
	sys *System
	fs  afero.Fs

	// This is the absolute path to the root of the project from which the
	// Generator will be operating.
	workDir string
}

//NewGenerator retrunes a new  Generator.
func NewGenerator(vm *otto.Otto, sys *System, fs afero.Fs) *Generator {
	return &Generator{
		vm:  vm,
		sys: sys,
		fs:  fs,
	}
}

//Build builds a project.
func (g *Generator) Build() error {
	steps := []struct {
		stage BuildStage
		exec  func() error
	}{
		{StageConfig, g.Config},
		{StageInit, g.Init},
		{StagePlan, g.plan},
		{StageExec, g.exec},
	}
	var err error
	for _, buildStep := range steps {
		err = buildStep.exec()
		if err != nil {
			return err
		}
	}
	return nil
}

//Init initializes the build process. Any stages after this will have the generator
//already bootstraped.
//
// It is possible to bootstrap the generator from the project( User's side) by
// providing an entry javascript file in the default path of
// scripts/init/index.js which will be executed and you can overide the default
// entry script which is evaluated internally
//
// Initialzation is offloaded to the javascript runtine of the generator..Any
// error returned is a build error.
func (g *Generator) Init() error {

	_ = g.vm.Set("sys", func(call otto.FunctionCall) otto.Value {
		data, err := json.Marshal(g.sys)
		if err != nil {
			panicOtto(err)
		}
		val, err := call.Otto.Call("JSON.parse", nil, string(data))
		if err != nil {
			panicOtto(err)
		}
		return val
	})
	_, err := g.vm.Eval(entryScript())
	if err != nil {
		return buildErr(StageInit, err.Error())
	}

	_ = g.vm.Set("fileTree", fileTree(g.fs, g.workDir))

	// evaluate project provided entry script if provided. We ignore if the file
	// is not provided but any errors arsing from evaluating a provided script is
	// a built error.
	entryFile := filepath.Join(scriptsDir, initDir, indexFile)
	err = g.evaluateFile(entryFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return buildErr(StageInit, err.Error())
		}
	}

	return nil
}

//fileTree provides an array of all files found in the root( which is suppose to
//be the working directory.
//
//The file index is built only once and evaluated once too, then it is cached.
//This operates in the otto runtime.
func fileTree(fs afero.Fs, root string) func(otto.FunctionCall) otto.Value {
	var tree []string
	var v otto.Value
	return func(call otto.FunctionCall) otto.Value {
		if tree != nil {
			return v
		}
		ferr := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				tree = append(tree, path)
			}
			return nil
		})
		if ferr != nil {
			tree = nil
			panicOtto(ferr.Error())
		}
		v, _ = call.Otto.ToValue(tree)
		return v
	}
}

func defaultSystem() *System {
	return &System{
		Boot: &Boot{
			ConfigiFile: configFile,
			PlanFile:    "index.js",
		},
	}
}

// opens the file in the specified path and evaluates it withing the context of
// the javascript runtine.
//
// The evaluated javascript code can mutate the global state. Use execFile to
// execute the javascript without mutating the state of the generato'r
// javascript runtime.
//
// TODO: (gernest) implement callFile if necessary
func (g *Generator) evaluateFile(path string) error {
	f, err := g.fs.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = g.vm.Eval(data)
	if err != nil {
		return err
	}
	return nil
}

func defaultVM(sys *System) *otto.Otto {
	return otto.New()
}

//Config configures the generator.
func (g *Generator) Config() error {
	if g.sys == nil {
		g.sys = defaultSystem()
	}
	if g.vm == nil {
		g.vm = defaultVM(g.sys)
	}
	// Properly set working directory/
	if g.workDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return buildErr(StageInit, err.Error())
		}
		g.workDir = wd
	}

	// ensure everything is relative to the working directory
	g.fs = afero.NewBasePathFs(g.fs, g.workDir)
	af := afero.Afero{Fs: g.fs}
	data, err := af.ReadFile(configFile)
	if err != nil {
		return buildErr(StageInit, err.Error())
	}
	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return buildErr(StageInit, err.Error())
	}
	g.sys.Config = cfg
	return nil
}

func (g *Generator) plan() error {
	return nil
}
func (g *Generator) exec() error {
	return nil
}
func (g *Generator) down() error {
	return nil
}
