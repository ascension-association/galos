package main

import (
	"fmt"
	"time"
	"context"

	execute "github.com/alexellis/go-execute/v2"
	"github.com/gokrazy/gokrazy"
)

var container = "docker.io/library/hello-world:latest"

func run(logging bool, exe string, args ...string) {
	cmd := execute.ExecTask{
		Command:     exe,
		Args:        args,
		StreamStdio: false,
	}

	res, err := cmd.Execute(context.Background())
	if err != nil {
		if logging {
			fmt.Errorf("Error: %v", err)
		}
	}

	if res.ExitCode != 0 {
		if logging {
			fmt.Errorf("Error: %v", res.Stderr)
		}
	}

	if logging {
		fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
	}
}

func main() {
	// create mount point
	run(false, "/usr/local/bin/busybox", "mkdir", "-p", "/perm/galos")

	// wait for network
	gokrazy.WaitForClock()

	// wait for containerd
	time.Sleep(3 * time.Second)

	// remove prior container, if applicable
	run(false, "/usr/local/bin/ctr", "task", "remove", "--force", "galos")
	run(false, "/usr/local/bin/ctr", "snapshot", "remove", "galos")
	run(false, "/usr/local/bin/ctr", "container", "remove", "galos")

	// pull container
	run(true, "/usr/local/bin/ctr", "image", "pull", container)

	// create container
	run(true, "/usr/local/bin/ctr", "container", "create", "--privileged", "--net-host", "--mount", "type=bind,src=/perm/galos,dst=/perm,options=rbind:rw", container, "galos")

	// run container
	run(true, "/usr/local/bin/ctr", "task", "start", "--detach", "galos")
}
