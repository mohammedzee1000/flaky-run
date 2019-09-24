package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mohammedzee1000/flaky-run/pkg/flakeexecutor"
)

func main() {
	var runCount, parallel int
	var logsDir string
	var cmd string
	var args []string
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, "run-logs")
	flag.IntVar(&runCount, "runs", 5, "Specify how many runs need to happen, defaults to 5")
	flag.IntVar(&parallel, "parallel", 1, "Specify parallel runs, defaults to 1")
	flag.StringVar(&logsDir, "logsdir", dir, "Specify where you want to store logs, defaults to currentdir/run-;ogs")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s [options] [command] [commandargs...]\n", os.Args[0], filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	ta := flag.Args()
	if len(ta) <= 0 {
		flag.Usage()
		os.Exit(0)
	}
	cmd = ta[0]
	if len(ta) > 1 {
		args = ta[1:]
	}
	err = flakeexecutor.NewCommandFlakeExecutor(runCount, parallel, logsDir, cmd, args...).Run()
	if err != nil {
		log.Fatal(err)
	}
}
