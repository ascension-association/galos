package main

import (
	"fmt"
	"time"
	"context"

	execute "github.com/alexellis/go-execute/v2"
	"github.com/gokrazy/gokrazy"
)

func run(exe string, args ...string) {
	cmd := execute.ExecTask{
		Command:     exe,
		Args:        args,
		StreamStdio: false,
	}

	res, err := cmd.Execute(context.Background())
	if err != nil {
		panic(err)
	}

	if res.ExitCode != 0 {
		panic("Non-zero exit code: " + res.Stderr)
	}

	fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}

func main() {
	// wait for network
	gokrazy.WaitForClock()

	// wait for containerd
	time.Sleep(3 * time.Second)

	run("/usr/local/bin/ctr", "version")
}
