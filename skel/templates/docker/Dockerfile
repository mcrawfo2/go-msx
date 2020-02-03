###
### Build container
###

FROM golang:1.12-stretch AS builder
ARG BUILD_FLAGS
ENV BUILDER_FLAGS -mod=vendor
ADD . /app/
WORKDIR /app
RUN make dist

###
### Distribution Container
###

FROM dockerhub.cisco.com/vms-platform-dev-docker/vms-base-stretch:3.8.0-370
ENV SERVER_STATICPATH /var/lib/${app.name}
EXPOSE ${server.port}
WORKDIR /var/run/${app.name}/
CMD ["/usr/bin/${app.name}"]
COPY --from=builder /app/dist/root/ /
