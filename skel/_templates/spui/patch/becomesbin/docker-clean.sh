#!/bin/bash
# Clean UI docker instances if exists

source $(pwd)/bin/vars.sh

# Required variables
SERVICEPACK_NAME="${SERVICEPACK_NAME:?}"
DOCKERFILE="${DOCKERFILE:?}"
FULL_VERSION="${FULL_VERSION:?}"

# Calculated variables
IMAGE_NAME="${SERVICEPACK_NAME}-ui"
IMAGE_VERSION="${FULL_VERSION}"

for stage_name in build test conformance dist; do
    echo "Removing $stage_name image"
    docker rmi -f "${IMAGE_NAME}_${stage_name}:${IMAGE_VERSION}"
done

docker rmi -f "dockerhub.cisco.com/vms-platform-dev-docker/${MSX_RELEASE}/latest/${IMAGE_NAME}:${IMAGE_VERSION}"
