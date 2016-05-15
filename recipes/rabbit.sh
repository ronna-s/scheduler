#!/bin/bash

apt-get install -y rabbitmq-server
rabbitmq-plugins enable rabbitmq_management

sudo service rabbitmq-server restart

if [ ! -f /usr/sbin/rabbitmqadmin ]; then
    cd /usr/sbin
    wget http://localhost:15672/cli/rabbitmqadmin
    chmod +x rabbitmqadmin
fi

sudo service rabbitmq-server restart
