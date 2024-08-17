# Container Images

This repo contains container images for various services we use in our projects.

Some highlights:

### [Cribl](https://cribl.io)

A hardened Cribl container image.

* Uses the [Distroless `cc-debian12`](https://github.com/GoogleContainerTools/distroless/blob/main/cc/README.md) base image.
* Runs as a non-root user. 
  * The `/opt/cribl` directory is owned by the non-root user allowing for some flexibility.
* `busybox` system utilities are installed.
* Runs the [s6-overlay](https://github.com/just-containers/s6-overlay) process supervisor.
* The Cribl process is wrapped using [tini](https://github.com/krallin/tini) for clean shutdown.
