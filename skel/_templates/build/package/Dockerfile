ARG BASE_IMAGE=dockerhub.cisco.com/vms-platform-dev-docker/msx-base-bullseye:latest

###
### Build container
###

FROM dockerhub.cisco.com/vms-platform-dev-docker/golang:1.16-bullseye AS builder
ARG BUILD_FLAGS
ENV BUILDER_FLAGS -mod=vendor
ADD . /app/
WORKDIR /app
RUN make dist

###
### Distribution Container
###

FROM ${BASE_IMAGE}
EXPOSE ${server.port}
WORKDIR /var/run/${app.name}/
CMD ["/usr/bin/${app.name}"]
COPY --from=builder /app/dist/root/ /
