package main

import (
	"fmt"
	"log"
	"context"

	execute "github.com/alexellis/go-execute/v2"
	"github.com/gokrazy/gokrazy"
)

var container = "docker.io/library/hello-world:latest"
//var executable = ""
//var arguments = ""

func run(logging bool, exe string, args ...string) {
	var cmd execute.ExecTask

	if logging {
		cmd = execute.ExecTask{
			Command:     exe,
			Args:        args,
			StreamStdio: true,
		}
	} else {
		cmd = execute.ExecTask{
			Command:     exe,
			Args:        args,
			StreamStdio: false,
			DisableStdioBuffer: true,
		}
	}

	res, err := cmd.Execute(context.Background())

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	if res.ExitCode != 0 {
		fmt.Errorf("Error: %v", res.Stderr)
	}
}

func main() {
	log.Println("Initializing network...")

	// wait for network
	gokrazy.WaitFor("net-online")

	log.Println("Initializing Galos...")

	// create mount point
	run(true, "/usr/local/bin/busybox", "mkdir", "-p", "/perm/galos")

	// remove prior container, if applicable
	run(false, "/usr/local/bin/ctr", "task", "remove", "--force", "galos")
	run(false, "/usr/local/bin/ctr", "snapshot", "remove", "galos")
	run(false, "/usr/local/bin/ctr", "container", "remove", "galos")

	// pull container
	log.Println("Pulling container image. Please wait, this may take a while...")
	run(false, "/usr/local/bin/ctr", "image", "pull", container)

	// create container
	log.Println("Creating container...")
	run(false, "/usr/local/bin/ctr", "container", "create", "--privileged", "--net-host", "--mount", "type=bind,src=/perm/galos,dst=/perm,options=rbind:rw", container, "galos")

	// run container
	log.Println("Running container...")
	run(true, "/usr/local/bin/ctr", "task", "start", "galos")
}
