#!/bin/bash -e

source $(pwd)/bin/vars.sh

# Required variables
SERVICEPACK_NAME="${SERVICEPACK_NAME:?}"
MSX_RELEASE="${MSX_RELEASE:?}"
TARBALL_NAME="${TARBALL_NAME:?}"
TARBALL_VERSION="${TARBALL_VERSION:?}"
ARTIFACTORY_USERNAME="${ARTIFACTORY_USERNAME:?}"
ARTIFACTORY_PASSWORD="${ARTIFACTORY_PASSWORD:?}"

# Constants
FILENAME="build/skyfall.tar"
DESTINATION_TARBALL_LEGACY="${SERVICEPACK_NAME}/${FULL_VERSION}/${TARBALL_NAME}-${TARBALL_VERSION}.tar"
DESTINATION_TARBALL="${MSX_RELEASE}/latest/${SERVICEPACK_NAME}/ui-${TARBALL_VERSION}.tar"
ARTIFACTORY_URL="https://engci-maven-master.cisco.com/artifactory/symphony-thirdparty/vms-3.0-binaries"
CREDS="${ARTIFACTORY_USERNAME}:${ARTIFACTORY_PASSWORD}"

echo "Transfer package to artifactory..."
# //TODO: pushing to the LEGACY destination to be removed when Harness is fully adopted.
curl -u ${CREDS} -X PUT ${ARTIFACTORY_URL}/${DESTINATION_TARBALL_LEGACY} -T ${FILENAME}

curl -u ${CREDS} -X PUT ${ARTIFACTORY_URL}/${DESTINATION_TARBALL} -T ${FILENAME}
echo "Done transfer."
