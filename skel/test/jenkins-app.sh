#!/bin/bash -xe

WORKSPACE=${WORKSPACE:?}
GOMSX=${WORKSPACE}/go-msx
GOBIN=${WORKSPACE}/go/bin
GEN=${GOMSX}/skel/test
SERVICE=${WORKSPACE}/dummyservice

export GOBIN

# Install skel locally
(
    cd ${GOMSX}
    mkdir -p "${GOBIN}"
    make install-skel
)

# Generate a microservice
(
    cd ${GEN}
    cat > generate.json <<EOF
{
    "generator": "app",
    "targetParent": "${WORKSPACE}",
    "appName": "dummyservice",
    "appDisplayName": "Dummy Microservice",
    "appDescription": "Verifies microservice generation",
    "serverPort": 9909,
    "serverContextPath": "/dummy",
    "appVersion": "3.10.0",
    "repository": "cockroach",
    "deploymentGroup": "dummyservice"
}
EOF

    rm -Rf ${SERVICE}
    ${GOBIN}/skel
)

# Generate some domains
(
    cd ${SERVICE}
    ${GOBIN}/skel generate-domain-system animal
    ${GOBIN}/skel generate-domain-tenant pet
    ${GOBIN}/skel generate-domain-tenant owner
    date > .random
)

# Deploy the jenkins job

# Force-push code to develop

