package main

// forked from https://gokrazy.org/packages/docker-containers/

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"context"
	"time"

	"github.com/gokrazy/gokrazy"
)

var container = "docker.io/library/hello-world:latest"

func ctr(args ...string) error {
	ctr := exec.Command("/usr/local/bin/ctr", args...)
	ctr.Env = expandPath(os.Environ())
	ctr.Stdin = os.Stdin
	ctr.Stdout = os.Stdout
	ctr.Stderr = os.Stderr
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	ctr = ctr.WithContext(ctx)
	if err := ctr.Run(); err != nil {
		return fmt.Errorf("%v: %v", ctr.Args, err)
	}
	return nil
}

func galos() error {
	// Ensure we have an up-to-date clock, which in turn also means that
	// networking is up.
	gokrazy.WaitForClock()
	
	// wait a few seconds for containerd to initialize
	time.Sleep(3 * time.Second)

	task, err := exec.Command("/usr/local/bin/ctr", "task", "list", "--quiet").Output()
	if err != nil {
		log.Print(err)
	}

	if strings.TrimRight(string(task), "\n") == "galos" {
		fmt.Printf("blah %v\n", "1")
	    /*
	    if err := ctr("task", "remove", "--force", "galos"); err != nil {
		    log.Print(err)
	    }

	    if err := ctr("snapshot", "remove", "galos"); err != nil {
		    log.Print(err)
	    }

	    if err := ctr("container", "remove", "galos"); err != nil {
		    log.Print(err)
	    }
	    */
	}
/*
	if err := ctr("image", "pull", container); err != nil {
		log.Print(err)
	}

	if err := ctr("container", "create",
		"--privileged", "--net-host",
		"--hostname", "galos",
		"--mount", "type=bind,src=/perm/galos,dst=/perm,options=rbind:rw",
		container, "galos"); err != nil {
		return err
	}

	if err := ctr("task", "start", "--detach", "galos"); err != nil {
		log.Print(err)
	}
*/
	return nil
}

func main() {
	makeDirectoryIfNotExists("/perm/galos")
	if err := galos(); err != nil {
		log.Fatal(err)
	}
}

// expandPath returns env, but with PATH= modified or added
// such that both /user and /usr/local/bin are included, which ctr needs.
func expandPath(env []string) []string {
	extra := "/user:/usr/local/bin"
	found := false
	for idx, val := range env {
		parts := strings.Split(val, "=")
		if len(parts) < 2 {
			continue // malformed entry
		}
		key := parts[0]
		if key != "PATH" {
			continue
		}
		val := strings.Join(parts[1:], "=")
		env[idx] = fmt.Sprintf("%s=%s:%s", key, extra, val)
		found = true
	}
	if !found {
		const busyboxDefaultPATH = "/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin"
		env = append(env, fmt.Sprintf("PATH=%s:%s", extra, busyboxDefaultPATH))
	}
	return env
}

// https://gist.github.com/ivanzoid/5040166bb3f0c82575b52c2ca5f5a60c
func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}
