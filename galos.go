// https://github.com/containerd/containerd/blob/main/docs/getting-started.md#implementing-your-own-containerd-client
package main

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/containerd/containerd/v2/pkg/cio"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/containerd/v2/pkg/namespaces"
)

var container = "docker.io/library/hello-world:latest"

func main() {
	if err := galos(); err != nil {
		log.Fatal(err)
	}
}

func galos() error {
	// remove prior running instance, if applicable
	taskList, err := exec.Command("/usr/local/bin/ctr", "--namespace galos", "task", "list", "--quiet").Output()
	if err != nil {
		log.Print(err)
	}
	if strings.TrimRight(string(taskList), "\n") == "galos" {
		if err := exec.Command("/usr/local/bin/ctr", "--namespace galos", "task", "remove", "--force", "galos").Run(); err != nil {
			log.Print(err)
		}
		if err := exec.Command("/usr/local/bin/ctr", "--namespace galos", "snapshot", "remove", "galos").Run(); err != nil {
			log.Print(err)
		}
		if err := exec.Command("/usr/local/bin/ctr", "--namespace galos", "container", "remove", "galos").Run(); err != nil {
			log.Print(err)
		}
	}

	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with namespace
	ctx := namespaces.WithNamespace(context.Background(), "galos")

	// pull the image
	image, err := client.Pull(ctx, container, containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	// create a container
	container, err := client.NewContainer(
		ctx,
		"galos",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("galos", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return err
	}
	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	log.Print(code)

	// call start on the task to execute the container
	if err := task.Start(ctx); err != nil {
		return err
	}

	return nil
}
