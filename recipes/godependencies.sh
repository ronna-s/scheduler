#!/bin/bash

export GOROOT=/opt/go
export GOPATH=/home/vagrant/gopath
sudo chown -R vagrant /home/vagrant/gopath --quiet
go get github.com/tools/godep
go get github.com/onsi/ginkgo/ginkgo 
go install github.com/onsi/ginkgo/ginkgo 
go get github.com/onsi/gomega
go get github.com/onsi/gomega
go get github.com/streadway/amqp