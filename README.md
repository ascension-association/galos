# Galos
Minimal OS to run a Podman container

### Overview
Inspired by projects like [CoreOS](https://github.com/coreos), [RancherOS](https://github.com/rancher/os), and [Talos Linux](https://github.com/siderolabs/talos), Galos is a minimal Linux OS powered by [gokrazy](https://github.com/gokrazy/gokrazy) that is primarily designed to run a single [Podman](https://github.com/containers/podman) container on IoT bare metal with as few executables as possible.

_Note: As you probably guessed, the name 'Galos' is a nod to **g**okrazy and T**alos** Linux (unless you're a Sidero Labs lawyer, in which case it's a tribute to the **Gal**ápag**os** Islands)_

### Instructions
1. Include in your gokrazy instance: `gok new`
   
2. Add Galos and its dependencies:

```
gok add github.com/gokrazy/iptables
gok add github.com/gokrazy/nsenter
gok add github.com/gokrazy/podman
gok add github.com/greenpau/cni-plugins/cmd/cni-nftables-portmap
gok add github.com/greenpau/cni-plugins/cmd/cni-nftables-firewall
gok add github.com/gokrazy/mkfs
gok add github.com/ascension-association/galos
```

3. Set desired container in your gokrazy config.json _PackageConfig_ section:

```
"github.com/ascension-association/galos": {
    "GoBuildFlags": [
        "-ldflags=-X main.container=docker.io/library/hello-world:latest"
    ]
}
```

4. If your target is AMD64, set `"GOARCH=amd64"` and add `"KernelPackage": "github.com/gokrazy/kernel.amd64",` 

5. If deploying via USB/SD at location /dev/sda: `gok overwrite --full /dev/sda` Otherwise, if you're targeting an already deployed instance: `gok update`

6. Verify it worked by going to the gokrazy dashboard, clicking on the `/user/galos` link and reviewing the logs

7. Optionally, once confirmed working, edit your gokrazy config.json again and remove the `"github.com/gokrazy/hello",`, `"github.com/gokrazy/fbstatus",` and `"github.com/gokrazy/mkfs",` packages (and even `"github.com/gokrazy/breakglass",` if you don't need SSH access), then run `gok update` to further minimize the device contents
