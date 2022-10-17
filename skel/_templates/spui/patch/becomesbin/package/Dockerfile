ARG BUILD_BASE=dockerhub.cisco.com/vms-platform-dev-docker/vms-cicd-jenkins-node:8
ARG DIST_BASE=nginx:1.17.4

### Build
FROM $BUILD_BASE AS build
ARG WORKSPACE=/workspace/ui
ADD . $WORKSPACE
WORKDIR $WORKSPACE
RUN ln -sfn $WORKSPACE /app
RUN env

# configure npm
ARG NPM_PROXY=http://proxy.esl.cisco.com:80
RUN npm config set proxy "$NPM_PROXY"
RUN npm config set https-proxy "$NPM_PROXY"
ARG NPM_VMS_REPO=http://engci-maven-master.cisco.com/artifactory/api/npm/vms-npm
RUN npm config set @nstehr:registry "$NPM_VMS_REPO"

# install dependencies
RUN npm install

# run the build
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build

### Tests
FROM build AS test
RUN npm run coverage

### Static Analysis
FROM build AS conformance
RUN npm install @nstehr/conformance-cli
RUN npm run conformance

### Distribution
FROM $DIST_BASE AS dist
RUN mkdir -p /usr/share/nginx/html/services/
COPY --from=build /app/build/* /usr/share/nginx/html/services/
