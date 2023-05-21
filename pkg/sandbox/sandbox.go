package sandbox

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type SandboxManager struct {
	logger *logrus.Logger
}

func NewSanboxManager(logger *logrus.Logger) *SandboxManager {
	return &SandboxManager{logger: logger}
}

func (sandboxManager *SandboxManager) CleanupIsolate(id int) {
	command := NewCommand(true).WithCleanup().Box(id).Command()
	cleanupCmd := exec.Command("isolate", command...)
	if err := cleanupCmd.Run(); err != nil {
		sandboxManager.logger.Fatal("unable to cleanup sanbox isolates", err)
	}
}

func (sandboxManager *SandboxManager) InitIsolate(id int, useControlGroups bool) (string, error) {
	sandboxManager.logger.Debug("initializing isolate ", id)
	sandboxManager.CleanupIsolate(id)

	command := NewCommand(useControlGroups).WithInit().Box(id).Command()
	cmd := exec.Command("isolate", command...)

	dir, err := cmd.Output()

	if err != nil {
		return "", errors.New(fmt.Sprintf("error while initializing sandbox isolate %s", err.Error()))
	}

	isolateDir := strings.TrimSpace(string(dir))
	return path.Join(isolateDir, "box"), nil
}
