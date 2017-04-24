#!/bin/bash

# Join the default ECS cluster
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

PATH=$PATH:/usr/local/bin

# Instance should be added to an security group that allows HTTP outbound
yum update -y

# Install packages
yum -y install jq nfs-utils python27 python27-pip

EC2_INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`

# NewRelic Infrastructure Agent
echo "license_key: ${newrelic_key}" | sudo tee -a /etc/newrelic-infra.yml
printf "[newrelic-infra]\nname=New Relic Infrastructure\nbaseurl=http://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64\nenable=1\ngpgcheck=0" | sudo tee -a /etc/yum.repos.d/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
sudo yum install newrelic-infra -y

# Using a different CIDR for Docker (not to conflict with coreIT VPC CIDR)
sed -i '$ d' /etc/sysconfig/docker
echo 'OPTIONS="--default-ulimit nofile=1024:4096 --bip=192.168.5.1/24"' >> /etc/sysconfig/docker
service docker stop
ip link set docker0 down
ip link delete docker0
service docker start

# Restarting ECS
start ecs
