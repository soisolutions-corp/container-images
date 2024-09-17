FROM --platform=$BUILDPLATFORM golang:1.23-bookworm AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/wait-all

FROM scratch
COPY --from=build /out/wait-all /usr/local/bin/wait-all
