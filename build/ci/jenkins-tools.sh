#!/bin/bash

GO_VERSION=1.20
GO_TARBALL=go${GO_VERSION}.linux-amd64.tar.gz

curl -LO https://dl.google.com/go/${GO_TARBALL}
sudo rm -rvf /usr/local/go/
sudo tar -C /usr/local -xzf ${GO_TARBALL}
sudo chown -R jenkins:jenkins /usr/local/go
rm -f ${GO_TARBALL} || true
