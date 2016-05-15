#!/bin/bash

tee /home/vagrant/.bash_profile <<EOF
alias l='ls -lh --color'
export GOROOT=/opt/go
export GOPATH=/home/vagrant/gopath
export PATH=\$PATH:\$GOPATH/bin
export SCHEDULER_HOME_DIR=\$GOPATH/src/github.com/ronna-s/scheduler
export PS1='\[\033[0;32m\]\u@\h \[\033[1;33m\]\w\[\033[0m\] '
EOF
