#!/bin/bash -e

source $(pwd)/bin/vars.sh

# Required variables
SERVICEPACK_NAME="${SERVICEPACK_NAME:?}"
DOCKERFILE="${DOCKERFILE:?}"
FULL_VERSION="${FULL_VERSION:?}"

# Optional variables
DOCKER_USERNAME=${DOCKER_USERNAME:-}
DOCKER_PASSWORD=${DOCKER_PASSWORD:-}

if [ "$DOCKER_USERNAME" != "" ] && [ "$DOCKER_PASSWORD" != "" ]; then
    docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}" dockerhub.cisco.com
fi

docker push "dockerhub.cisco.com/vms-platform-dev-docker/${MSX_RELEASE}/latest/${IMAGE_NAME}:${IMAGE_VERSION}"
