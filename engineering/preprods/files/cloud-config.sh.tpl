#!/bin/bash

# Join the default ECS cluster
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

PATH=$PATH:/usr/local/bin

# Instance should be added to an security group that allows HTTP outbound
yum update -y

# Install packages
yum -y install jq nfs-utils python27 python27-pip

EC2_INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`

#NewRelic Installation
rpm -Uvh https://download.newrelic.com/pub/newrelic/el5/i386/newrelic-repo-5-3.noarch.rpm
yum install newrelic-sysmond -y
usermod -a -G docker newrelic
nrsysmond-config --set license_key=${newrelic_key}
sed -i '/#hostname*/c\hostname=\"${newrelic_hostname} ('$EC2_INSTANCE_ID')\"' /etc/newrelic/nrsysmond.cfg
sed -i '/#labels*/c\labels=${newrelic_labels}' /etc/newrelic/nrsysmond.cfg
/etc/init.d/newrelic-sysmond start

# Using a different CIDR for Docker (not to conflict with coreIT VPC CIDR)
sed -i '$ d' /etc/sysconfig/docker
echo 'OPTIONS="--default-ulimit nofile=1024:4096 --bip=192.168.5.1/24"' >> /etc/sysconfig/docker
service docker stop
ip link set docker0 down
ip link delete docker0
service docker start

# Restarting ECS
start ecs
