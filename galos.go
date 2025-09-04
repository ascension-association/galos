package main

import (
	"fmt"
	"context"
	"log"

	execute "github.com/alexellis/go-execute/v2"
	"github.com/gokrazy/gokrazy"
)

var container = "docker.io/library/hello-world:latest"
//var executable = ""
//var arguments = ""

func run(logging bool, exe string, args ...string) {
	cmd := execute.ExecTask{
		Command:     exe,
		Args:        args,
		StreamStdio: false,
	}

	res, err := cmd.Execute(context.Background())
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	if res.ExitCode != 0 {
		fmt.Errorf("Error: %v", res.Stderr)
	}

	if logging {
		fmt.Printf(res.Stdout)
	}
}

func main() {
	log.Print("Initializing...")

	// create mount point
	run(false, "/usr/local/bin/busybox", "mkdir", "-p", "/perm/galos")

	// wait for network
	gokrazy.WaitFor("net-online")

	// remove prior container, if applicable
	run(false, "/usr/local/bin/ctr", "task", "remove", "--force", "galos")
	run(false, "/usr/local/bin/ctr", "snapshot", "remove", "galos")
	run(false, "/usr/local/bin/ctr", "container", "remove", "galos")

	// pull container
	fmt.Printf("Pulling container image. Please wait, this may take a while...")
	run(false, "/usr/local/bin/ctr", "image", "pull", container)

	// create container
	fmt.Printf("Creating container...")
	run(false, "/usr/local/bin/ctr", "container", "create", "--privileged", "--net-host", "--mount", "type=bind,src=/perm/galos,dst=/perm,options=rbind:rw", container, "galos")

	// run container
	fmt.Printf("Running container...")
	run(true, "/usr/local/bin/ctr", "task", "start", "galos")

	// show results
	fmt.Printf("Done:")
	run(true, "/usr/local/bin/ctr", "task", "exec", "--exec-id", "oneshot", "galos", "ps aux | head -n 2")
}
