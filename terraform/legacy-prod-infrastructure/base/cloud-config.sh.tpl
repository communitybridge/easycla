#!/bin/bash

# Join the default ECS cluster
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

# Adding gelf & syslog logging drivers
echo 'ECS_AVAILABLE_LOGGING_DRIVERS=["json-file","syslog","awslogs","gelf"]' >> /etc/ecs/ecs.config

PATH=$PATH:/usr/local/bin

# Instance should be added to an security group that allows HTTP outbound
yum update -y

# Install packages
yum -y install jq nfs-utils python27 python27-pip

# Get region of EC2 from instance metadata
EC2_AVAIL_ZONE=`curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone`
EC2_REGION="`echo \"$EC2_AVAIL_ZONE\" | sed -e 's:\([0-9][0-9]*\)[a-z]*\$:\\1:'`"
EC2_INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`
EC2_PRIVATE_IP=`curl -s http://169.254.169.254/latest/meta-data/local-ipv4`

# NewRelic Infrastructure Agent
echo "license_key: ${newrelic_key}" | sudo tee -a /etc/newrelic-infra.yml
printf "[newrelic-infra]\nname=New Relic Infrastructure\nbaseurl=http://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64\nenable=1\ngpgcheck=0" | sudo tee -a /etc/yum.repos.d/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
sudo yum install newrelic-infra -y

# Restarting ECS
start ecs

# Install FileBeat Agent
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-5.4.0-x86_64.rpm
sudo rpm -vi filebeat-5.4.0-x86_64.rpm
cat > '/etc/filebeat/filebeat.yml' <<-EOF
filebeat.prospectors:

- input_type: log

  paths:
    - /var/log/cloud-init*.log
    - /var/log/ecs/*.log
    - /var/log/messages
    - /var/log/secure
    - /var/log/yum.log
    - /var/log/cron
    - /var/log/maillog
    - /var/log/docker
    - /var/log/lastlog
    - /var/log/audit/*.log

#================================ General =====================================

name: infrastructure-ecs-cluster
tags: ["system"]
fields:
  sys_name: infrastructure
  sys_env: production
  sys_region: ${region}

#================================ Outputs =====================================
output.logstash:
  hosts: ["${region}.logstash.service.consul:5044"]
EOF

service filebeat start

initctl restart newrelic-infra