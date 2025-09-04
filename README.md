# Galos
Minimal OS to run a containerd container

### Overview
Inspired by projects like [CoreOS](https://github.com/coreos), [RancherOS](https://github.com/rancher/os), and [Talos Linux](https://github.com/siderolabs/talos), Galos is a minimal Linux OS powered by [gokrazy](https://github.com/gokrazy/gokrazy) that is primarily designed to run a single [containerd](https://github.com/containerd/containerd) container on IoT bare metal with as few executables as possible.

_Note: As you probably guessed, the name 'Galos' is a nod to **G**oKrazy and T**alos** Linux (unless you're a Sidero Labs lawyer, in which case it's a tribute to the **Gal**apag**os** Islands)_

### Prerequisites
1. Install Vim

2. Install Go >= 1.24

3. Run `go install github.com/gokrazy/tools/cmd/gok@main`

### Instructions
1. Create a new gokrazy instance: `gok new`

2. Run `gok edit` and add your desired container in the _PackageConfig_ section:

```
"github.com/ascension-association/galos": {
    "GoBuildFlags": [
        "-ldflags=-X main.container=ghcr.io/apptainer/lolcow:latest"
    ]
}
```

If the container doesn't have an automatic entrypoint command or you want to run your own, use this format:

```
"github.com/ascension-association/galos": {
    "GoBuildFlags": [
        "-ldflags=-X main.container=ghcr.io/void-linux/void-musl-busybox:latest -X 'main.task=cat /etc/os-release'"
    ]
}
```

_Important: the `ctr` command-line tool, which Galos uses, requires **fully-qualified image references**, including the **registry domain and the tag**, such as `docker.io/library/nginx:latest`, and cannot be abbreviated to just `docker.io/library/nginx` or `nginx` or `nginx:latest`._

_Note: if no container is provided, Galos defaults to `docker.io/library/hello-world:latest`_

3. **IF running on x86/amd64**, do the following then save:

  - add this line under the "Hostname" line: `"KernelPackage": "github.com/gokrazy/kernel.amd64",`
  - change `"GOARCH=arm64"` to `"GOARCH=amd64"`

4. Add Galos and its dependencies:

```
gok add github.com/gokrazy/mkfs
gok add github.com/ascension-association/containerd
gok add github.com/ascension-association/galos
```

5. If deploying via USB/SD at location /dev/sda: `gok overwrite --full /dev/sda` Otherwise, if you're targeting an already deployed instance: `gok update`

IF deploying via USB/SD, plug into target device and boot from it. Use the URL provided in the output of the `gok overwrite` step to load in your source machine's browser (note: you may need to replace 'hello' with the IP address of the target device).

6. Verify it worked by going to the gokrazy dashboard, clicking on the `/user/galos` link and reviewing the logs

7. Optionally, once confirmed working, edit your gokrazy config.json again and remove the `"github.com/gokrazy/hello",`, `"github.com/gokrazy/fbstatus",` and `"github.com/gokrazy/mkfs",` packages (and even `"github.com/gokrazy/breakglass",` if you don't need SSH access), then run `gok update` to further minimize the device contents

