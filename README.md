# Galos
Minimal OS to run a Podman container

### Overview
Inspired by projects like [CoreOS](https://github.com/coreos), [RancherOS](https://github.com/rancher/os), and [Talos Linux](https://github.com/siderolabs/talos), Galos is a minimal Linux OS powered by [gokrazy](https://github.com/gokrazy/gokrazy) that is primarily designed to run a single [Podman](https://github.com/containers/podman) container on IoT bare metal with as few executables as possible.

_Note: As you probably guessed, the name 'Galos' is a nod to **g**okrazy and T**alos** Linux (unless you're a Sidero Labs lawyer, in which case it's a tribute to the **Gal**ápag**os** Islands)_

### Instructions
1. Include in your gokrazy instance: `gok add github.com/ascension-association/galos`
2. Set desired container in your gokrazy config.json PackageConfig section:

```
"github.com/ascension-association/galos": {
    "GoBuildFlags": [
        "-ldflags=-X main.container=docker.io/library/hello-world:latest"
    ]
}
```

3. Deploy: `gok update`

