#!/bin/bash -e
#Build UI and UI Unit Tests

source $(pwd)/bin/vars.sh

# Required variables
SERVICEPACK_NAME="${SERVICEPACK_NAME:?}"
DOCKERFILE="${DOCKERFILE:?}"
FULL_VERSION="${FULL_VERSION:?}"
IMAGE_NAME="${IMAGE_NAME:?}"
IMAGE_VERSION="${IMAGE_VERSION:?}"

# Optional variables
BUILD_IMAGE=${BUILD_IMAGE:-dockerhub.cisco.com/docker.io/node:12}
DIST_IMAGE=${DIST_IMAGE:-dockerhub.cisco.com/docker.io/nginx:1.17.4}
DOCKER_USERNAME=${DOCKER_USERNAME:-}
DOCKER_PASSWORD=${DOCKER_PASSWORD:-}
SKIP_CONFORMANCE=${SKIP_CONFORMANCE:-true}
HTTPS_PROXY=${HTTPS_PROXY:-}
NO_PROXY=${NO_PROXY:-}
WORKSPACE=${WORKSPACE:-}

# Calculated variables
SOURCE_PATH="$(pwd)"

echo "Build ${SERVICEPACK_NAME} UI"

stage() {
    local stage_name=${1:?}
    local skip=${2:-false}

    if [ "$skip" == "true" ]; then
        echo "Skipping stage '$stage_name'"
        return 0
    fi

    echo "Executing stage '$stage_name'"

    export DOCKER_BUILDKIT=1

    if [ "$WORKSPACE" == "" ]; then
        PROGRESS=""
    else
        PROGRESS="--progress=plain"
    fi

    set +e
    docker build ${PROGRESS} \
        --build-arg WORKSPACE="$SOURCE_PATH" \
        --build-arg NPM_PROXY="$HTTPS_PROXY" \
        --build-arg NO_PROXY \
        --build-arg BUILD_BASE="$BUILD_IMAGE" \
        --build-arg DIST_BASE="$DIST_IMAGE" \
        --target "${stage_name}" \
        -f "${DOCKERFILE}" \
        -t "${IMAGE_NAME}_${stage_name}:${IMAGE_VERSION}" \
        ${SOURCE_PATH}

    result=$?
    set -e

    if [ "$result" != "0" ]; then
        echo "Failed stage '$stage_name'"
        return ${result}
    fi

    echo "Completed stage '$stage_name'"
    return 0
}

extract() {
    local stage_name=${1:?}
    local inside=${2:?}
    local outside=${3:?}
    local image_name="${IMAGE_NAME}_${stage_name}:${IMAGE_VERSION}"

    echo "Extracting ${inside} to ${outside}"

    mkdir -p "${outside}"
    CONTAINER_ID=$(docker create ${image_name})
    docker cp "${CONTAINER_ID}:${inside}" "${outside}"
    docker rm -f ${CONTAINER_ID}
}

tag() {
    local stage_name=${1:?}
    docker tag "${IMAGE_NAME}_${stage_name}:${IMAGE_VERSION}" "dockerhub.cisco.com/vms-platform-dev-docker/${MSX_RELEASE}/latest/${IMAGE_NAME}:${IMAGE_VERSION}"
}

docker builder prune -a -f

docker pull "${BUILD_IMAGE}"
docker pull "${DIST_IMAGE}"

stage build
extract build "$SOURCE_PATH/build/." "build/"

stage test
extract test "$SOURCE_PATH/test/." "test/"

stage conformance "${SKIP_CONFORMANCE}"

stage dist
tag dist
