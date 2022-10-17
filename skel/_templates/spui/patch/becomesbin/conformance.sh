#!/bin/bash -e

source $(pwd)/bin/vars.sh

# Required variables
SOURCE_PATH="${SOURCE_PATH:?}"

# Optional variables
SKIP_CONFORMANCE=${SKIP_CONFORMANCE:-false}

if [ "$SKIP_CONFORMANCE" == "true" ]; then
    exit 0
fi

if [ -d "$SOURCE_PATH/frontend" ]; then
    SOURCE_PATH="${SOURCE_PATH}/frontend"
fi

node --harmony node_modules/@nstehr/conformance-cli/bin/index.js -d "${SOURCE_PATH}" -p
