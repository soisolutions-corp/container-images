FROM debian:bookworm-slim as git

WORKDIR /tmp

# renovate: datasource=github-tags depName=git/git
ARG GIT_VERSION=2.46.2

# Install dependencies
RUN apt-get update && apt-get install -y build-essential gettext autoconf libz-dev libssl-dev

# Download from github mirror
RUN curl -L https://github.com/git/git/archive/refs/tags/v${GIT_VERSION}.tar.gz | tar zxvf -

# Compile git as a static binary
WORKDIR /tmp/git-${GIT_VERSION}
RUN make configure && ./configure prefix=/usr/local/ CFLAGS="${CFLAGS} -static"
RUN make && make install

# Export the binary
FROM scratch
COPY --from=git /usr/local/bin/git /usr/local/bin/git
