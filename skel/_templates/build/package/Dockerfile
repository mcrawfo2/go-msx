ARG BASE_IMAGE=dockerhub.cisco.com/vms-platform-dev-docker/msx-base-bullseye:latest

###
### Build container
###

FROM dockerhub.cisco.com/docker.io/golang:1.18-bullseye AS builder
ARG BUILD_FLAGS
ENV BUILDER_FLAGS -mod=vendor
ADD . /app/
WORKDIR /app
RUN make dist
RUN cp build/package/docker-entrypoint.sh dist/root/
RUN chmod 755 dist/root/docker-entrypoint.sh

###
### Distribution Container
###

FROM ${BASE_IMAGE}
EXPOSE ${server.port}
ENV SERVICE_BIN "/usr/bin/${app.name}"
WORKDIR /var/run/${app.name}/
CMD ["/usr/bin/${app.name}"]
COPY --from=builder /app/dist/root/ /
ENTRYPOINT [ "/docker-entrypoint.sh" ]
