package engine

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/uvite/gvm/pkg/js"
	"github.com/uvite/gvm/pkg/lib"
	"github.com/uvite/gvm/pkg/loader"
	"github.com/uvite/gvm/pkg/metrics"
	"gopkg.in/guregu/null.v3"
	"net/url"
)

func (gvm *Gvm)  GetSimpleRunner(filename, data string, opts ...interface{}) (*js.Runner, error) {
	var (
		rtOpts = lib.RuntimeOptions{CompatibilityMode: null.NewString("base", true)}

		fsResolvers = map[string]afero.Fs{"file": afero.NewMemMapFs(), "https": afero.NewMemMapFs()}
	)
	for _, o := range opts {
		switch opt := o.(type) {
		case afero.Fs:
			fsResolvers["file"] = opt
		case map[string]afero.Fs:
			fsResolvers = opt
		case lib.RuntimeOptions:
			rtOpts = opt
		//case *logrus.Logger:
		//	logger = opt
		default:
			fmt.Printf("unknown test option %q", opt)
		}
	}
	logger := logrus.New()
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	return js.New(
		&lib.TestPreInitState{
			Logger:         logger,
			RuntimeOptions: rtOpts,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		},
		&loader.SourceData{
			URL:  &url.URL{Path: filename, Scheme: "file"},
			Data: []byte(data),
		},
		fsResolvers,
	)
}
