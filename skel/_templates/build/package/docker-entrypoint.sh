#!/bin/bash

# Supported ENV variables:
# * POPULATE: Unset, or one of the following to populate options: "database", "all", "resourceString", "serviceMetadata", "serviceCatalog", "customRolesAndCapabilities", "deviceAction", "secretPolicy", "billingScript"
# Note: currently no support for mode

if [ -z "$POPULATE" ]; then
  exec $SERVICE_BIN
fi

MIGRATE_COMMAND="$SERVICE_BIN migrate"

if [ "$POPULATE" = "database" ]; then
    exec $MIGRATE_COMMAND
fi

if [ "$POPULATE" = "all" ]; then
    eval $MIGRATE_COMMAND
    checkResult=$(echo $?)
    if [ $checkResult != 0 ]; then
        echo "Failed to execute migrate as part of populate all!"
        exit 1
    fi
fi

exec $SERVICE_BIN populate $POPULATE
