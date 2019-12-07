###
### Build container
###

FROM golang:1.12-stretch AS builder
ADD . /app/
WORKDIR /app
RUN make dist

###
### Distribution Container
###

FROM dockerhub.cisco.com/vms-platform-dev-docker/vms-base-stretch:3.7.0-339
COPY --from=builder /app/dist/root/ /
EXPOSE ${server.port}
WORKDIR /var/run/${app.name}/
CMD ["/usr/bin/${app.name}"]
