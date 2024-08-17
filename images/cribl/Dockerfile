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

# Tini
FROM base AS tini
# renovate: datasource=github-tags depName=krallin/tini
ARG TINI_VERSION=v0.19.0
ARG TARGETARCH
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static-${TARGETARCH} /tini
RUN chmod +x /tini

# Statically compiled git
FROM ghcr.io/soisolutions-corp/git:2.44.0 AS git

FROM busybox:1.36.1-uclibc as busybox

FROM cribl/cribl:4.8.1 as cribl

# Cribl requires glibc, so we use the cc-debian12 image
# https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md
FROM gcr.io/distroless/cc-debian12 as final
COPY --from=busybox /bin /bin
COPY --from=git /usr/local/bin/git /usr/local/bin/git
COPY --from=cribl --chown=65532:65532 /opt/cribl /opt/cribl
COPY --from=tini /tini /tini
COPY --from=s6-overlay /s6/ /

# Make /var/run a symlink to /run so we can run this in an read-only container
RUN /bin/rm -rf /var/run
RUN /bin/ln -s /run /var/run

# Copy in all of our container files
COPY container-files /

# Set user to `nonroot`
USER 65532

# Flatten the final image
FROM scratch
COPY --from=final / /

# Make s6-overlay less chatty
ENV S6_VERBOSITY=0

ENTRYPOINT ["/init"]
