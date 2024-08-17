# Git

Compiles `git` as a static binary. You can use this binary in other container images that require `git`.

This image is meant to included in another multi-stage image build. 

## Usage

```Dockerfile
FROM ghcr.io/soisolutions-corp/git:<semver> AS git

# ... your container image build steps here

COPY --from=git /usr/local/bin/git /usr/local/bin/git
```
