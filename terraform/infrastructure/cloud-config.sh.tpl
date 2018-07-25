#!/bin/bash

echo vm.max_map_count=262144 >> /etc/sysctl.conf
sysctl -w vm.max_map_count=262144

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
EC2_INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`

# NewRelic Infrastructure Agent
echo "license_key: ${newrelic_key}" | sudo tee -a /etc/newrelic-infra.yml
printf "[newrelic-infra]\nname=New Relic Infrastructure\nbaseurl=http://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64\nenable=1\ngpgcheck=0" | sudo tee -a /etc/yum.repos.d/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
sudo yum install newrelic-infra -y

service docker start

# Restarting ECS
start ecs

# Install FileBeat Agent
#curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-5.4.0-x86_64.rpm
#sudo rpm -vi filebeat-5.4.0-x86_64.rpm
#cat > '/etc/filebeat/filebeat.yml' <<-EOF
#filebeat.prospectors:
#
#- input_type: log
#
#  paths:
#    - /var/log/cloud-init*.log
#    - /var/log/ecs/*.log
#    - /var/log/messages
#    - /var/log/secure
#    - /var/log/yum.log
#    - /var/log/cron
#    - /var/log/maillog
#    - /var/log/docker
#    - /var/log/lastlog
#    - /var/log/audit/*.log
#
#================================ General =====================================
#
#name: staging-cluster
#tags: ["system"]
#fields:
#  sys_name: staging-ecs-cluster
#  sys_env: staging
#  sys_region: ${aws_region}
#
#================================ Outputs =====================================
#output.logstash:
#  hosts: ["${aws_region}.logstash.service.consul:5044"]
#EOF
#
#service filebeat start

initctl restart newrelic-infra

# Create mount point
mkdir /mnt/storage
mkdir /mnt/storage/consul
mkdir /mnt/storage/nexus
mkdir /mnt/storage/mongodb

CONSUL_DIR_SRC=$EC2_AVAIL_ZONE.${consul_efs}.efs.$EC2_REGION.amazonaws.com:/
NEXUS_DIR_SRC=$EC2_AVAIL_ZONE.${nexus_efs}.efs.$EC2_REGION.amazonaws.com:/
MONGODB_DIR_SRC=$EC2_AVAIL_ZONE.${mongodb_efs}.efs.$EC2_REGION.amazonaws.com:/

# Mount EFS file system
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $CONSUL_DIR_SRC /mnt/storage/consul
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $NEXUS_DIR_SRC /mnt/storage/nexus
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $MONGODB_DIR_SRC /mnt/storage/mongodb

#Backup fstab
cp -p /etc/fstab /etc/fstab.back-$(date +%F)

#Append line to fstab
echo -e "$CONSUL_DIR_SRC \t\t /mnt/storage/consul \t\t nfs \t\t defaults \t\t 0 \t\t 0" | tee -a /etc/fstab
echo -e "$NEXUS_DIR_SRC \t\t /mnt/storage/nexus \t\t nfs \t\t defaults \t\t 0 \t\t 0" | tee -a /etc/fstab
echo -e "$MONGODB_DIR_SRC \t\t /mnt/storage/mongodb \t\t nfs \t\t defaults \t\t 0 \t\t 0" | tee -a /etc/fstab