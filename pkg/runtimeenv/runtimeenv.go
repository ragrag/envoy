package runtimeenv

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ragrag/envoy/pkg/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _   = runtime.Caller(0)
	runtimesRoot = filepath.Join(filepath.Dir(b), "../../runtimes")
)

type Runtime struct {
	ID                string `yaml:"id"`
	File              string `yaml:"file"`
	Dir               string
	CompileScriptPath string
	RunScriptPath     string
	Meta              struct {
		Language string `yaml:"language"`
		Version  string `yaml:"version"`
	} `yaml:"meta"`
}

type RuntimeProvider struct {
	logger   *logrus.Logger
	runtimes map[string]*Runtime
}

func NewRuntimeProvider(logger *logrus.Logger) *RuntimeProvider {
	return &RuntimeProvider{logger: logger, runtimes: make(map[string]*Runtime)}
}

func (runtimeProvider *RuntimeProvider) Load() error {
	runtimeProvider.logger.Info("loading runtimes")

	err := filepath.Walk(runtimesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New("error while parsing runtimes")
		}

		if info.IsDir() && filepath.Base(path) == "runtime" {
			return filepath.SkipDir
		}

		if !info.IsDir() && filepath.Base(path) == "runtime.yaml" {
			parsedRuntime, err := parseRuntimeYaml(path)

			if err != nil {
				return errors.New("error while parsing runtimes")
			}

			parsedRuntime.Dir = filepath.Dir(path)
			compileScriptExists, err := util.FileExists(filepath.Join(parsedRuntime.Dir, "compile.sh"))
			if err != nil {
				return errors.New("error while parsing runtimes")
			}

			if compileScriptExists {
				parsedRuntime.CompileScriptPath = filepath.Join(parsedRuntime.Dir, "compile.sh")
			} else {
				parsedRuntime.CompileScriptPath = ""
			}

			parsedRuntime.RunScriptPath = filepath.Join(parsedRuntime.Dir, "run.sh")

			runtimeProvider.runtimes[parsedRuntime.ID] = parsedRuntime
		}

		return nil
	})

	runtimeProvider.logger.Infof("%d runtimes loaded", len(runtimeProvider.runtimes))

	return err
}

func (runtimeProvider *RuntimeProvider) GetRuntime(id string) (*Runtime, error) {
	runtime, ok := runtimeProvider.runtimes[id]

	if !ok {
		return nil, errors.New("Runtime is not found")
	}

	return runtime, nil
}

func (runtimeProvider *RuntimeProvider) GetAvailableRuntimes() []*Runtime {
	return maps.Values(runtimeProvider.runtimes)
}

func parseRuntimeYaml(path string) (*Runtime, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &Runtime{}, err
	}

	var parsedRuntime Runtime
	err = yaml.Unmarshal(data, &parsedRuntime)
	if err != nil {
		return &Runtime{}, err
	}

	return &parsedRuntime, nil
}
