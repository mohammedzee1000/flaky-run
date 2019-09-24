package runner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mohammedzee1000/flaky-run/pkg/utils"
	"github.com/pkg/errors"
)

//CommandRunner runs a specific command
type CommandRunner struct {
	command string
	args    []string
}

//NewCommandRunner creates a new command runner
func NewCommandRunner(command string, args ...string) *CommandRunner {
	return &CommandRunner{
		command: command,
		args:    args,
	}
}

func (r *CommandRunner) writeError(logsDir string, err error) error {
	msg := "no errors"
	if err != nil {
		msg = err.Error()
	}
	return utils.WriteFile(filepath.Join(logsDir, "error.log"), []byte(msg))
}

func (r *CommandRunner) writeLogs(logsDir string, out string) error {
	return utils.WriteFile(filepath.Join(logsDir, "out.log"), []byte(out))
}

func (r *CommandRunner) String() string {
	return fmt.Sprint(r.command, " ", strings.Join(r.args, " "))
}

//Run runs the command
func (r *CommandRunner) Run(logsDir string) (bool, error) {
	var success bool
	o, err := utils.RunCommand(r.command, r.args...)
	if err == nil {
		success = true
	}
	err = r.writeError(logsDir, err)
	if err != nil {
		return false, errors.Wrap(err, "failed to write error messages")
	}
	err = r.writeLogs(logsDir, o)
	if err != nil {
		return false, errors.Wrap(err, "failed to write output logs")
	}
	return success, nil
}
