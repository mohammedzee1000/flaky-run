package flakeexecutor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/mohammedzee1000/flaky-run/pkg/runner"
	"github.com/mohammedzee1000/flaky-run/pkg/utils"
	"github.com/pkg/errors"
)

//CommandFlakeExecutor executes the Flake
type CommandFlakeExecutor struct {
	runner   *runner.CommandRunner
	runCount int
	parallel int
	logsDir  string
}

//NewCommandFlakeExecutor returns a new flake runner
func NewCommandFlakeExecutor(runCount int, parallel int, logsDir string, command string, args ...string) *CommandFlakeExecutor {
	return &CommandFlakeExecutor{
		runner:   runner.NewCommandRunner(command, args...),
		runCount: runCount,
		parallel: parallel,
		logsDir:  logsDir,
	}
}

func (e *CommandFlakeExecutor) init() error {
	err := os.RemoveAll(e.logsDir)
	if err != nil {
		return errors.Wrap(err, "failed to delete existing logs dir")
	}
	err = os.MkdirAll(filepath.Join(e.logsDir, "runs"), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to create logdir %s", e.logsDir)
	}
	return nil
}

func (e *CommandFlakeExecutor) runForFlake() (int, int) {
	var success, fail int
	var wg sync.WaitGroup
	var tr int
	var counter int
	status := make(chan bool, e.runCount)
	counter = e.runCount
	for {
		if counter <= 0 {
			break
		}
		if counter < e.parallel {
			tr = counter
		} else {
			tr = e.parallel
		}
		for i := 0; i < tr; i++ {
			wg.Add(1)
			logsDir := filepath.Join(e.logsDir, "runs", strconv.Itoa(counter))
			err := os.MkdirAll(logsDir, os.ModePerm)
			if err != nil {
				log.Fatal(errors.Wrapf(err, "failed to create runs %d", counter))
			}
			go func(counter int, logsDir string) {
				s, err := e.runner.Run(logsDir)
				if err != nil {
					log.Fatal(err)
				}
				status <- s
				wg.Done()
			}(counter, logsDir)
			counter = counter - 1
		}
		wg.Wait()
	}
	close(status)
	for ele := range status {
		if ele {
			success = success + 1
		} else {
			fail = fail + 1
		}
	}
	return success, fail
}

func (e *CommandFlakeExecutor) summerize(success int, fail int) error {
	sumFile := filepath.Join(e.logsDir, "summary.log")
	var data string
	data = fmt.Sprintf("COMMAND=%s\nSUCCESS=%d\nFAILS=%d", e.runner.String(), success, fail)
	err := utils.WriteFile(sumFile, []byte(data))
	if err != nil {
		return errors.Wrap(err, "failed to summerise")
	}
	fmt.Println(data)
	return nil
}

//Run run validation for flake
func (e *CommandFlakeExecutor) Run() error {
	if e.parallel <= 0 {
		return errors.New("parallel runners must be a natuarl number")
	}
	err := e.init()
	if err != nil {
		return err
	}
	s, f := e.runForFlake()
	err = e.summerize(s, f)
	if err != nil {
		return err
	}
	return nil
}
