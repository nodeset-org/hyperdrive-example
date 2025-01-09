# The builder for building the service
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
    cd /hyperdrive-example/service && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/hd-service

# The daemon image
FROM debian:bookworm-slim AS service
COPY --from=builder /build/hd-service /usr/bin/hd-service

# Container entry point
ENTRYPOINT [ "/usr/bin/hd-service", "-c", "/hd/config/service-cfg.yaml", "-k", "/hd/secret" ]