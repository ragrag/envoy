package sandbox

import (
	"fmt"
	"strconv"
)

type CommandBuilder struct {
	cg      bool
	command []string
}

func NewCommand(cg bool) *CommandBuilder {
	return &CommandBuilder{cg: cg}
}

func (command *CommandBuilder) WithRun() *CommandBuilder {
	command.command = append(command.command, "--run")
	if command.cg {
		command.command = append(command.command, "--cg")
	}
	return command
}

func (command *CommandBuilder) WithCleanup() *CommandBuilder {
	command.command = append(command.command, "--cleanup")
	if command.cg {
		command.command = append(command.command, "--cg")
	}
	return command
}

func (command *CommandBuilder) WithInit() *CommandBuilder {
	command.command = append(command.command, "--init")
	if command.cg {
		command.command = append(command.command, "--cg")
	}
	return command
}

func (command *CommandBuilder) Box(id int) *CommandBuilder {
	command.command = append(command.command, "-b", strconv.Itoa(id))
	return command
}

func (command *CommandBuilder) Processes() *CommandBuilder {
	command.command = append(command.command, "-p")
	return command
}

func (command *CommandBuilder) BinPathEnv() *CommandBuilder {
	pathEnv := `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin`
	command.command = append(command.command, "-E", pathEnv)
	return command
}

func (command *CommandBuilder) BoxRootEnv(dir string) *CommandBuilder {
	boxRootEnv := fmt.Sprintf(`BOX_ROOT=%s`, dir)
	command.command = append(command.command, "-E", boxRootEnv)
	return command
}

func (command *CommandBuilder) DirRW(dir string) *CommandBuilder {
	command.command = append(command.command, "-d", fmt.Sprintf("%s:%s", dir, "rw"))
	return command
}

func (command *CommandBuilder) DirNoExec(dir string) *CommandBuilder {
	command.command = append(command.command, "-d", fmt.Sprintf("%s:%s", dir, "noexec"))
	return command
}

func (command *CommandBuilder) Bash(script string) *CommandBuilder {
	command.command = append(command.command, "--", "/bin/bash", script)
	return command
}

func (command *CommandBuilder) Stdin(file string) *CommandBuilder {
	command.command = append(command.command, "-i", file)
	return command
}

func (command *CommandBuilder) Stdout(file string) *CommandBuilder {
	command.command = append(command.command, "-o", file)
	return command
}

func (command *CommandBuilder) Stderr(file string) *CommandBuilder {
	command.command = append(command.command, "-r", file)
	return command
}

func (command *CommandBuilder) StderrToStdout() *CommandBuilder {
	command.command = append(command.command, "--stderr-to-stdout")
	return command
}

func (command *CommandBuilder) TimeLimit(seconds float64) *CommandBuilder {
	command.command = append(command.command, fmt.Sprintf("--time=%f", seconds))
	command.command = append(command.command, fmt.Sprintf("--wall-time=%f", (seconds*2)+1))
	return command
}

func (command *CommandBuilder) MemoryLimit(kilobytes int) *CommandBuilder {
	if command.cg {
		command.command = append(command.command, fmt.Sprintf("--cg-mem=%d", kilobytes))
	} else {
		command.command = append(command.command, fmt.Sprintf("--mem=%d", kilobytes))
	}
	return command
}

func (command *CommandBuilder) Meta(file string) *CommandBuilder {
	command.command = append(command.command, "-M", file)
	return command
}

func (command *CommandBuilder) Command() []string {
	return command.command
}
