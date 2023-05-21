package engine

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/ragrag/envoi/pkg/infra"
	"github.com/ragrag/envoi/pkg/runtimeenv"
	"github.com/ragrag/envoi/pkg/sandbox"
	"github.com/ragrag/envoi/pkg/util"

	"github.com/sirupsen/logrus"
)

type Engine struct {
	logger          *logrus.Logger
	config          *infra.EngineConfig
	runtimeProvider *runtimeenv.RuntimeProvider
	sandboxManager  *sandbox.SandboxManager
	isolatePool     chan *isolateWorker
}

type isolateWorker struct {
	id  int
	dir string
}

type isolateMeta struct {
	time          float64
	memory        int
	oomKilled     bool
	timedOut      bool
	internalError bool
	message       string
}

type ResultStatus string

const (
	COMPILATION_ERROR     ResultStatus = "COMPILATION_ERROR"
	RUNTIME_ERROR         ResultStatus = "RUNTIME_ERROR"
	TIME_LIMIT_EXCEEDED   ResultStatus = "TIME_LIMIT_EXCEEDED"
	MEMORY_LIMIT_EXCEEDED ResultStatus = "MEMORY_LIMIT_EXCEEDED"
	WRONG_ANSWER          ResultStatus = "WRONG_ANSWER"
	SUCCESS               ResultStatus = "SUCCESS"
)

type RunOptions struct {
	TimeLimit   *float64
	MemoryLimit *int
}

type RunParams struct {
	Language string
	Code     string
	Stdin    string
	Options  *RunOptions
}

type RunResult struct {
	Output string       `json:"output"`
	Status ResultStatus `json:"status"`
	Time   float64      `json:"time,omitempty"`
	Memory int          `json:"memory,omitempty"`
}

type TestCase struct {
	Input          string
	ExpectedOutput string
}

type JudgeParams struct {
	Language  string
	Code      string
	TestCases []TestCase
	Options   *RunOptions
}

type TestCaseResult struct {
	Output         string       `json:"output,omitempty"`
	ExpectedOutput string       `json:"expectedOutput,omitempty"`
	Input          string       `json:"input,omitempty"`
	Status         ResultStatus `json:"status"`
	Time           float64      `json:"time,omitempty"`
	Memory         int          `json:"memory,omitempty"`
}

type JudgeVerdict struct {
	Output         string       `json:"output,omitempty"`
	ExpectedOutput string       `json:"expectedOutput,omitempty"`
	Input          string       `json:"input,omitempty"`
	Status         ResultStatus `json:"status"`
	TotalTime      float64      `json:"totalTime,omitempty"`
	TotalMemory    int          `json:"totalMemory,omitempty"`
}

type JudgeResult struct {
	Verdict *JudgeVerdict     `json:"verdict"`
	Results []*TestCaseResult `json:"results"`
}

func NewEngine(logger *logrus.Logger, config *infra.Config, runtimeProvider *runtimeenv.RuntimeProvider, sandboxManager *sandbox.SandboxManager) *Engine {
	return &Engine{logger: logger, config: &config.Engine, runtimeProvider: runtimeProvider, sandboxManager: sandboxManager}
}

func (engine *Engine) Ignite() error {
	engine.isolatePool = make(chan *isolateWorker, engine.config.WorkerCount)

	for i := 0; i < engine.config.WorkerCount; i++ {
		isolateDir, err := engine.sandboxManager.InitIsolate(i, engine.config.UseControlGroups)
		if err != nil {
			return err
		}

		err = os.Chmod(isolateDir, 0777)
		if err != nil {
			return err
		}

		engine.isolatePool <- &isolateWorker{id: i, dir: isolateDir}
	}

	return nil
}

func (engine *Engine) RunCode(params *RunParams) (*RunResult, error) {
	isolate := <-engine.isolatePool
	defer engine.freeIsolate(isolate)

	engine.logger.Debug("isolate worker pooled ", isolate.id)

	err := isolate.cleanup()
	if err != nil {
		return nil, err
	}

	runtime, err := engine.runtimeProvider.GetRuntime(params.Language)
	if err != nil {
		return nil, err
	}

	err = engine.copyFilesToIsolate(runtime, isolate, params.Code)
	if err != nil {
		return nil, err
	}

	if runtime.CompileScriptPath != "" {
		errResult, err := engine.compile(runtime, isolate)
		if errResult != nil {
			return errResult.(*RunResult), nil
		}
		if err != nil {
			return nil, err
		}
	}

	return engine.run(runtime, isolate, params.Stdin, params.Options)
}

func (engine *Engine) JudgeCode(params *JudgeParams) (*JudgeResult, error) {
	isolate := <-engine.isolatePool
	defer engine.freeIsolate(isolate)

	engine.logger.Debug("isolate worker pooled ", isolate.id)

	err := isolate.cleanup()
	if err != nil {
		return nil, err
	}

	runtime, err := engine.runtimeProvider.GetRuntime(params.Language)
	if err != nil {
		return nil, err
	}

	err = engine.copyFilesToIsolate(runtime, isolate, params.Code)
	if err != nil {
		return nil, err
	}

	if runtime.CompileScriptPath != "" {
		errResult, err := engine.compile(runtime, isolate)
		if errResult != nil {
			verdict := &JudgeVerdict{Output: errResult.(*RunResult).Output, Status: errResult.(*RunResult).Status}
			results := make([]*TestCaseResult, 0)

			return &JudgeResult{Results: results, Verdict: verdict}, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return engine.judge(runtime, isolate, params.TestCases, params.Options)
}

func (engine *Engine) compile(runtime *runtimeenv.Runtime, isolate *isolateWorker) (interface{}, error) {
	engine.logger.Debug("compilation started")

	command := sandbox.NewCommand(engine.config.UseControlGroups).
		WithRun().
		Box(isolate.id).
		Processes().
		BinPathEnv().
		BoxRootEnv(isolate.dir).
		DirRW(isolate.dir).
		DirNoExec("/etc").
		StderrToStdout().
		Bash("compile.sh").
		Command()

	cmd := exec.Command("isolate", command...)

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	err := cmd.Run()

	if err != nil {
		stdout := stdoutBuf.String()

		if len(stdout) > 0 {
			return &RunResult{
				Output: stdout,
				Status: COMPILATION_ERROR,
			}, nil
		}

		return nil, errors.New(err.Error())
	}

	return nil, nil
}

func (engine *Engine) run(runtime *runtimeenv.Runtime, isolate *isolateWorker, stdin string, options *RunOptions) (*RunResult, error) {
	engine.logger.Debug("running binary")

	timeLimit := engine.config.TimeLimitSeconds
	memoryLimit := engine.config.MemoryLimitKB

	if options != nil && options.MemoryLimit != nil {
		memoryLimit = *options.MemoryLimit
	}
	if options != nil && options.TimeLimit != nil {
		timeLimit = *options.TimeLimit
	}

	command := sandbox.NewCommand(engine.config.UseControlGroups).
		WithRun().
		Box(isolate.id).
		Processes().
		BinPathEnv().
		BoxRootEnv(isolate.dir).
		DirRW(isolate.dir).
		DirNoExec("/etc").
		TimeLimit(timeLimit).
		MemoryLimit(memoryLimit).
		StderrToStdout().
		Meta(path.Join(isolate.dir, "meta.txt"))

	if len(stdin) > 0 {
		err := ioutil.WriteFile(path.Join(isolate.dir, "stdin.txt"), []byte(stdin), 0644)
		if err != nil {
			return nil, fmt.Errorf("error while preparing runtime files: %s", err.Error())
		}
		command = command.Stdin("stdin.txt")
	}

	command = command.Bash("run.sh")

	cmd := exec.Command("isolate", command.Command()...)

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	err := cmd.Run()
	stdout := stdoutBuf.String()

	if err != nil {
		meta, err := parseMeta(isolate, engine.config.UseControlGroups)
		if err != nil {
			return nil, err
		}

		if meta.oomKilled {
			return &RunResult{
				Output: stdout,
				Status: MEMORY_LIMIT_EXCEEDED,
				Time:   meta.time,
				Memory: meta.memory,
			}, nil
		}

		if meta.timedOut {
			return &RunResult{
				Output: stdout,
				Status: TIME_LIMIT_EXCEEDED,
				Time:   meta.time,
				Memory: meta.memory,
			}, nil
		}

		if meta.internalError {
			return nil, errors.New("internal error during binary execution")
		}

		return &RunResult{
			Output: stdout,
			Status: RUNTIME_ERROR,
			Time:   meta.time,
			Memory: meta.memory,
		}, nil
	}

	meta, err := parseMeta(isolate, engine.config.UseControlGroups)
	if err != nil {
		return nil, err
	}

	return &RunResult{
		Output: stdout,
		Status: SUCCESS,
		Time:   meta.time,
		Memory: meta.memory,
	}, nil
}

func (engine *Engine) judge(runtime *runtimeenv.Runtime, isolate *isolateWorker, testCases []TestCase, options *RunOptions) (*JudgeResult, error) {
	totalTime := 0.0
	totalMemory := 0

	var results []*TestCaseResult
	var verdict *JudgeVerdict

	for i := 0; i < len(testCases); i++ {
		err := ioutil.WriteFile(path.Join(isolate.dir, "stdin.txt"), []byte(testCases[i].Input), 0644)
		if err != nil {
			return nil, fmt.Errorf("error while preparing runtime files: %s", err.Error())
		}

		timeLimit := engine.config.TimeLimitSeconds
		memoryLimit := engine.config.MemoryLimitKB

		if options != nil && options.MemoryLimit != nil {
			memoryLimit = *options.MemoryLimit
		}
		if options != nil && options.TimeLimit != nil {
			timeLimit = *options.TimeLimit
		}

		command := sandbox.NewCommand(engine.config.UseControlGroups).
			WithRun().
			Box(isolate.id).
			Processes().
			BinPathEnv().
			BoxRootEnv(isolate.dir).
			DirRW(isolate.dir).
			DirNoExec("/etc").
			TimeLimit(timeLimit).
			MemoryLimit(memoryLimit).
			Stdin("stdin.txt").
			StderrToStdout().
			Meta(path.Join(isolate.dir, "meta.txt")).
			Bash("run.sh").
			Command()

		cmd := exec.Command("isolate", command...)

		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf

		err = cmd.Run()
		stdout := stdoutBuf.String()

		var testCaseResult *TestCaseResult

		if err != nil {
			meta, err := parseMeta(isolate, engine.config.UseControlGroups)
			if err != nil {
				return nil, err
			}

			if meta.oomKilled {
				testCaseResult = &TestCaseResult{
					Output:         stdout,
					ExpectedOutput: testCases[i].ExpectedOutput,
					Input:          testCases[i].Input,
					Status:         MEMORY_LIMIT_EXCEEDED,
					Time:           meta.time,
					Memory:         meta.memory,
				}

				results = append(results, testCaseResult)
			} else if meta.timedOut {
				testCaseResult = &TestCaseResult{
					Output:         stdout,
					ExpectedOutput: testCases[i].ExpectedOutput,
					Input:          testCases[i].Input,
					Status:         TIME_LIMIT_EXCEEDED,
					Time:           meta.time,
					Memory:         meta.memory,
				}

				results = append(results, testCaseResult)
			} else if meta.internalError {
				return nil, errors.New("internal error during binary execution")
			} else {
				testCaseResult = &TestCaseResult{
					Output:         stdout,
					ExpectedOutput: testCases[i].ExpectedOutput,
					Input:          testCases[i].Input,
					Status:         RUNTIME_ERROR,
					Time:           meta.time,
					Memory:         meta.memory,
				}

				results = append(results, testCaseResult)
			}
		}

		meta, err := parseMeta(isolate, engine.config.UseControlGroups)
		if err != nil {
			return nil, err
		}

		totalTime += meta.time
		totalMemory += meta.memory

		if stdout != testCases[i].ExpectedOutput {
			testCaseResult = &TestCaseResult{
				Output:         stdout,
				ExpectedOutput: testCases[i].ExpectedOutput,
				Input:          testCases[i].Input,
				Status:         WRONG_ANSWER,
				Time:           meta.time,
				Memory:         meta.memory,
			}

			results = append(results, testCaseResult)
		} else {
			results = append(results, &TestCaseResult{
				Output:         stdout,
				ExpectedOutput: testCases[i].ExpectedOutput,
				Input:          testCases[i].Input,
				Status:         SUCCESS,
				Time:           meta.time,
				Memory:         meta.memory,
			})
		}

		if testCaseResult != nil && verdict == nil {
			verdict = &JudgeVerdict{
				Output:         testCaseResult.Output,
				ExpectedOutput: testCaseResult.ExpectedOutput,
				Input:          testCaseResult.Input,
				Status:         testCaseResult.Status,
			}
		}
	}

	if verdict != nil {
		verdict.TotalTime = totalTime
		verdict.TotalMemory = totalMemory
		return &JudgeResult{
			Verdict: verdict,
			Results: results,
		}, nil
	}

	return &JudgeResult{
		Verdict: &JudgeVerdict{
			Status:      SUCCESS,
			TotalTime:   totalTime,
			TotalMemory: totalMemory,
		},
		Results: results,
	}, nil
}

func (engine *Engine) copyFilesToIsolate(runtime *runtimeenv.Runtime, isolate *isolateWorker, code string) error {
	engine.logger.Debug("copying runtime files to isolate", isolate.dir)

	var wg sync.WaitGroup
	wg.Add(3)

	copyErr := "error while preparing runtime files: %s"
	errChan := make(chan error, 3)

	// compile.sh
	go func() {
		defer wg.Done()
		if runtime.CompileScriptPath != "" {
			copyCompileScript := exec.Command("cp", runtime.CompileScriptPath, isolate.dir)
			err := copyCompileScript.Run()
			if err != nil {
				errChan <- fmt.Errorf(copyErr, err.Error())
			}
		}
	}()

	// run.sh
	go func() {
		defer wg.Done()
		copyRunScript := exec.Command("cp", runtime.RunScriptPath, isolate.dir)
		err := copyRunScript.Run()
		if err != nil {
			errChan <- fmt.Errorf(copyErr, err.Error())
		}
	}()

	// source code
	go func() {
		defer wg.Done()
		err := ioutil.WriteFile(path.Join(isolate.dir, runtime.File), []byte(code), 0777)
		if err != nil {
			errChan <- fmt.Errorf(copyErr, err.Error())
		}
	}()

	wg.Wait()
	close(errChan)

	err := <-errChan
	if err != nil {
		return err
	}

	return nil
}

// can be buffered to parse line-by-line and exit early to squeeze some perf gains (albeit very negligible)
func parseMeta(isolate *isolateWorker, useControlGroups bool) (*isolateMeta, error) {
	meta, err := ioutil.ReadFile(path.Join(isolate.dir, "meta.txt"))
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]string)

	for _, line := range strings.Split(string(meta), "\n") {
		if len(line) > 0 {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return nil, errors.New("internal error during runtime parsing")
			}
			parsed[parts[0]] = parts[1]
		}
	}

	timeStr := parsed["time"]
	time, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return nil, errors.New("internal error during runtime parsing")
	}

	var memKey string
	if useControlGroups {
		memKey = "cg-mem"
	} else {
		memKey = "max-rss"
	}

	memStr := parsed[memKey]

	mem, err := strconv.Atoi(memStr)
	if err != nil {
		return nil, errors.New("error while parsing metadata")
	}

	oomKilledVal, _ := strconv.Atoi(parsed["cg-oom-killed"])
	oomKilled := oomKilledVal == 1

	message := parsed["message"]
	status := parsed["status"]

	timedOut := status == "TO"
	internalError := status == "XX"

	return &isolateMeta{time: time, memory: mem, oomKilled: oomKilled, timedOut: timedOut, message: message, internalError: internalError}, nil
}

func (isolate *isolateWorker) cleanup() error {
	err := util.RemoveContents(isolate.dir)
	return err
}

func (engine *Engine) freeIsolate(worker *isolateWorker) {
	engine.isolatePool <- worker
}
