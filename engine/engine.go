package engine

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/uvite/gvm/pkg/js"
	"github.com/uvite/gvm/pkg/lib"
	"github.com/uvite/gvm/pkg/lib/consts"
	"github.com/uvite/gvm/pkg/loader"
	"github.com/uvite/gvm/pkg/metrics"
	"os"
)

type Gvm struct {
	Ctx context.Context
	*goja.Runtime
	Vu     lib.ActiveVU
	Logger *logrus.Logger
	Runner *js.Runner
}

func NewGvm() (*Gvm, error) {
	gvm := Gvm{}
	return &gvm, nil
}

func (gvm *Gvm) Load(file string) error {

	fs := afero.NewOsFs()
	pwd, _ := os.Getwd()
	logger := logrus.New()
	gvm.Logger = logger
	code, err := loader.ReadSource(logger, fmt.Sprintf("%s/%s", pwd, file), pwd, map[string]afero.Fs{"file": fs}, nil)
	if err != nil {
		return fmt.Errorf("couldn't load file: %s", err)
	}
	//fmt.Println(code.Data)

	rtOpts := lib.RuntimeOptions{}
	r, err := gvm.GetSimpleRunner("/script.js", fmt.Sprintf(`
			//import {Nats} from 'k6/x/nats';
			//import ta from 'k6/x/ta';
			//import {sleep} from 'k6'; 

			%s

			`, code.Data),
		fs, rtOpts)

	//	fmt.Println(err)

	gvm.Runner = r
	gvm.Runtime = r.Bundle.Vm
	if err != nil {
		return fmt.Errorf("couldn't set exported options with merged values: %w", err)

	}
	ch := make(chan metrics.SampleContainer, 100)
	//ctx, _ := context.WithCancel(context.Background())
	//defer cancel()
	err = r.Setup(gvm.Ctx, ch)

	initVU, err := r.NewVU(gvm.Ctx, 1, 1, ch)

	vu := initVU.Activate(&lib.VUActivationParams{RunContext: gvm.Ctx})
	gvm.Vu = vu

	return nil
}
func (gvm *Gvm) Run() (goja.Value, error) {

	gvm.Vu.RunOnce()
	v, ok := gvm.Vu.RunDefault()
	if ok != nil {
		return nil, ok
	}
	return v, nil
}
func (gvm *Gvm) Set(name string, value any) {
	err := gvm.Runner.Bundle.Vm.Set(name, value)
	if err != nil {
		fmt.Println(err)
	}

}

func (gvm *Gvm) ExecFunc(fun string) (goja.Value, error) {
	r := gvm.Runner
	if !r.IsExecutable(fun) {
		// do not init a new transient VU or execute setup() if it wasn't
		// actually defined and exported in the script
		gvm.Logger.Debugf("%s() is not defined or not exported, skipping!", fun)
		return nil, nil
	}
	gvm.Logger.Debugf("Running %s()...", consts.SetupFn)

	setupCtx, setupCancel := context.WithTimeout(gvm.Ctx, r.GetTimeoutFor(consts.SetupFn))
	defer setupCancel()
	out := make(chan metrics.SampleContainer, 100)
	v, err := r.RunPart(setupCtx, out, fun, nil)
	if err != nil {
		return nil, err
	}
	// r.setupData = nil is special it means undefined from this moment forward
	if goja.IsUndefined(v) {

		return nil, nil
	}
	return v, err
}
