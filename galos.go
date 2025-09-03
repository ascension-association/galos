package main

import (
	"fmt"
	"time"

	execute "github.com/alexellis/go-execute/v2"
	"context"
)

func main() {
	// wait a few seconds for containerd to initialize
	time.Sleep(3 * time.Second)

	cmd := execute.ExecTask{
		Command:     "/usr/local/bin/ctr",
		Args:        []string{"version"},
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
