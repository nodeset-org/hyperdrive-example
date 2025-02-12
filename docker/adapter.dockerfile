# The builder for building the adapter
FROM --platform=${BUILDPLATFORM} golang:1.22-bookworm AS builder
ARG TARGETOS TARGETARCH BUILDPLATFORM
COPY . /hyperdrive-example
RUN if [ "$BUILDPLATFORM" = "linux/amd64" -a "$TARGETARCH" = "arm64" ]; then \
    # Install the GCC cross compiler
    apt update && apt install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu && \
    export CC=aarch64-linux-gnu-gcc && export CC_FOR_TARGET=gcc-aarch64-linux-gnu; \
    elif [ "$BUILDPLATFORM" = "linux/arm64" -a "$TARGETARCH" = "amd64" ]; then \
    apt update && apt install -y gcc-x86-64-linux-gnu g++-x86-64-linux-gnu && \
    export CC=x86_64-linux-gnu-gcc && export CC_FOR_TARGET=gcc-x86-64-linux-gnu; \
    fi && \
    cd /hyperdrive-example/adapter && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/hd-adapter

# The daemon image
FROM debian:bookworm-slim AS adapter
COPY --from=builder /build/hd-adapter /usr/bin/hd-adapter

# Install docker compose - TEMPORARY for beta
RUN apt update && \
    apt install -y ca-certificates curl && \
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc && \
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt update && \
    apt install -y docker-ce-cli docker-compose-plugin

# Container entry point
ENTRYPOINT [ "/usr/bin/hd-adapter" ]
CMD [ "idle" ]