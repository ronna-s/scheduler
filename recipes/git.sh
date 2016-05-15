#!/bin/bash

apt-get install -y git-cvs
echo 'ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $*' > /home/vagrant/ssh
chmod +x /home/vagrant/ssh
