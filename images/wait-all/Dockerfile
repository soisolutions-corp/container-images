FROM --platform=$BUILDPLATFORM golang:1.22-bookworm AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/wait-all

FROM busybox:1.36.1-uclibc as busybox

FROM scratch
COPY --from=busybox /bin /bin
COPY --from=build /out/wait-all /usr/local/bin/wait-all
