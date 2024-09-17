# Cribl

This is a custom Cribl container image. The Cribl binary is obtained from the official Cribl container image, but we have focused on the following:

* Security
  * Based on the [Distroless `cc-debian12`](https://github.com/GoogleContainerTools/distroless/blob/main/cc/README.md) for enhanced security.
  * Runs as a non-root user, but the `/opt/cribl` directory is owned by this non-root user permitting some flexibility.
* Extensibility
  * Utilizes the [s6-overlay](https://github.com/just-containers/s6-overlay) process supervisor instead of a BASH entrypoint script.
* Performance
  * All dependencies are statically compiled.
  * The Cribl process is wrapped using a [custom init handler](../../images/wait-all) for clean shutdown.
* Size
  * As of version 4.8, the official Cribl container image is 460MB. This image is 330MB.

## Notes

In order to shut down cleanly, we have also created an empty binary file to "trick" Cribl into thinking that it is running under Systemd so that the s6-overlay can correctly take the process fully down and back up. This is necessary because Cribl would otherwise fork a new process and end up with an unmanaged zombie.

## Variants

`gitless` - We provide a container image that does not include the git binary. This is useful for worker nodes because they do not require git.

