#!/bin/bash

# golang
go_version="1.4.2"
go_binary="/opt/go/bin/go"
if [ ! -f "${go_binary}" ]; then
    tar_basename="go${go_version}.linux-amd64.tar.gz"
    cd /opt
    wget --no-verbose "https://storage.googleapis.com/golang/${tar_basename}"
    tar xzf "${tar_basename}"
    ln -s "${go_binary}" /bin/go
fi

export GOROOT=/opt/go
export GOPATH=/home/vagrant/gopath

