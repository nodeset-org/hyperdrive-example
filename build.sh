#!/bin/bash

# Print a failure message to stderr and exit
fail() {
    MESSAGE=$1
    RED='\033[0;31m'
    RESET='\033[;0m'
    >&2 echo -e "\n${RED}**ERROR**\n$MESSAGE${RESET}\n"
    exit 1
}


# Builds the adapter image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_adapter() {
    echo "Building adapter image..."
    # If uploading, make and push a manifest
    if [ "$UPLOAD" = true ]; then
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-example-adapter:$VERSION -f docker/adapter.dockerfile --push . || fail "Error building adapter image."
    elif [ "$LOCAL_UPLOAD" = true ]; then
        if [ -z "$LOCAL_DOCKER_REGISTRY" ]; then
            fail "LOCAL_DOCKER_REGISTRY must be set to upload to a local registry."
        fi
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t $LOCAL_DOCKER_REGISTRY/nodeset/hyperdrive-example-adapter:$VERSION -f docker/adapter.dockerfile --push . || fail "Error building adapter image."
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-example-adapter:$VERSION -f docker/adapter.dockerfile . || fail "Error building adapter image."
    fi
    echo "done!"
}


# Builds the service image and pushes it to Docker Hub
# NOTE: You must install qemu first; e.g. sudo apt-get install -y qemu qemu-user-static
build_service() {
    echo "Building service image..."
    # If uploading, make and push a manifest
    if [ "$UPLOAD" = true ]; then
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-example-service:$VERSION -f docker/service.dockerfile --push . || fail "Error building service image."
    elif [ "$LOCAL_UPLOAD" = true ]; then
        if [ -z "$LOCAL_DOCKER_REGISTRY" ]; then
            fail "LOCAL_DOCKER_REGISTRY must be set to upload to a local registry."
        fi
        docker buildx build --rm --platform=linux/amd64,linux/arm64 --build-arg BINARIES_PATH=build/$VERSION -t $LOCAL_DOCKER_REGISTRY/nodeset/hyperdrive-example-service:$VERSION -f docker/service.dockerfile --push . || fail "Error building service image."
    else
        docker buildx build --rm --load --build-arg BINARIES_PATH=build/$VERSION -t nodeset/hyperdrive-example-service:$VERSION -f docker/service.dockerfile . || fail "Error building service image."
    fi
    echo "done!"
}


# Builds the module package
build_package() {
    echo -n "Building module package... "
    tar cfJ build/$VERSION/hyperdrive-example.zip package/* || fail "Error building module package."
    echo "done!"
}


# Tags the 'latest' Docker Hub image
tag_latest() {
    echo -n "Tagging 'latest' Docker images... "
    docker tag nodeset/hyperdrive-example-adapter:$VERSION nodeset/hyperdrive-example-adapter:latest
    docker tag nodeset/hyperdrive-example-service:$VERSION nodeset/hyperdrive-example-service:latest
    echo "done!"

    if [ "$UPLOAD" = true ]; then
        echo -n "Pushing to Docker Hub... "
        docker push nodeset/hyperdrive-example-adapter:latest
        docker push nodeset/hyperdrive-example-service:latest
        echo "done!"
    else
        echo "The image tags only exist locally. Rerun with -u to upload to Docker Hub."
    fi
}


# =================
# === Main Body ===
# =================

# Parse arguments
while getopts "adspoluv:" FLAG; do
    case "$FLAG" in
        a) ADAPTER=true SERVICE=true PACKAGE=true ;;
        d) ADAPTER=true ;;
        s) SERVICE=true ;;
        p) PACKAGE=true ;;
        o) LOCAL_UPLOAD=true ;;
        l) LATEST=true ;;
        u) UPLOAD=true ;;
        v) VERSION="$OPTARG" ;;
        *) usage ;;
    esac
done
if [ -z "$VERSION" ]; then
    usage
fi

# Cleanup old artifacts
rm -rf build/$VERSION/*
mkdir -p build/$VERSION

# Make a multiarch builder, ignore if it's already there
docker buildx create --name multiarch-builder --driver docker-container --use > /dev/null 2>&1

# Build the artifacts
if [ "$ADAPTER" = true ]; then
    build_adapter
fi
if [ "$SERVICE" = true ]; then
    build_service
fi
if [ "$PACKAGE" = true ]; then
    build_package
fi
if [ "$LATEST" = true ]; then
    tag_latest
fi