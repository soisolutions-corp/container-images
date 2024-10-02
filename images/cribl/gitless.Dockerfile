FROM debian:bookworm-slim AS base

# S6 Overlay
FROM base AS s6-install-base
RUN apt-get -qq update && apt-get install -y wget xz-utils

FROM s6-install-base AS s6-arch-arm64
ARG S6_ARCH=aarch64

FROM s6-install-base AS s6-arch-amd64
ARG S6_ARCH=x86_64

FROM s6-arch-${TARGETARCH} AS s6-overlay
# renovate: datasource=github-tags depName=just-containers/s6-overlay
ARG S6_OVERLAY_VERSION=v3.2.0.0
WORKDIR /s6
RUN wget -q -O - https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz | tar Jxpf - -C /s6
RUN wget -q -O - https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-${S6_ARCH}.tar.xz | tar Jxpf - -C /s6

# Custom process supervisor
FROM ghcr.io/soisolutions-corp/wait-all:1.0.1 as wait-all

FROM busybox:1.37.0-uclibc as busybox

FROM cribl/cribl:4.8.2 as cribl

# Cribl requires glibc, so we use the cc-debian12 image
# https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md
FROM gcr.io/distroless/cc-debian12 as final
COPY --from=busybox /bin /bin
COPY --from=cribl --chown=65532:65532 /opt/cribl /opt/cribl
COPY --from=wait-all /usr/local/bin/wait-all /usr/local/bin/wait-all
COPY --from=s6-overlay /s6/ /

# Make /var/run a symlink to /run so we can run this in an read-only container
RUN /bin/rm -rf /var/run
RUN /bin/ln -s /run /var/run

# Copy in all of our container files
COPY container-files /

# Flatten the final image
FROM scratch
COPY --from=final / /

# Set user to `nonroot`
USER 65532

# Make s6-overlay less chatty
ENV S6_VERBOSITY=0

ENTRYPOINT ["/init"]
