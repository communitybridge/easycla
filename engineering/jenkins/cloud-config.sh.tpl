#!/bin/bash

yum update -y
yum install -y gcc java-1.8.0-openjdk git yum-utils openssl-devel install jq nfs-utils python27 python27-pip

curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-5.3.1-x86_64.rpm
rpm -vi filebeat-5.3.1-x86_64.rpm
echo "license_key: ${newrelic_license}" | tee -a /etc/newrelic-infra.yml
printf "[newrelic-infra]\nname=New Relic Infrastructure\nbaseurl=http://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64\nenable=1\ngpgcheck=0" | tee -a /etc/yum.repos.d/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
yum install newrelic-infra -y

wget -O /etc/yum.repos.d/jenkins.repo http://pkg.jenkins-ci.org/redhat-stable/jenkins.repo
rpm --import http://pkg.jenkins-ci.org/redhat-stable/jenkins-ci.org.key
yum install jenkins -y

# Get region of EC2 from instance metadata
EC2_AVAIL_ZONE=`curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone`
EC2_REGION="`echo \"$EC2_AVAIL_ZONE\" | sed -e 's:\([0-9][0-9]*\)[a-z]*\$:\\1:'`"

# Create mount point
mkdir /var/lib/jenkins

# Get EFS FileSystemID attribute
#Instance needs to be added to a EC2 role that give the instance at least read access to EFS
EFS_FILE_SYSTEM_ID=${efs_id}

# Instance needs to be a member of security group that allows 2049 inbound/outbound
#The security group that the instance belongs to has to be added to EFS file system configuration
#Create variables for source and target
DIR_SRC=$EC2_AVAIL_ZONE.$EFS_FILE_SYSTEM_ID.efs.$EC2_REGION.amazonaws.com:/
DIR_TGT=/var/lib/jenkins

# Mount EFS file system
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $DIR_SRC $DIR_TGT

#Backup fstab
cp -p /etc/fstab /etc/fstab.back-$(date +%F)

#Append line to fstab
echo -e "$DIR_SRC \t\t $DIR_TGT \t\t nfs \t\t defaults \t\t 0 \t\t 0" | tee -a /etc/fstab

sleep 230
mount -a
service jenkins start
chkconfig jenkins on

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
    - /var/log/jenkins/*.log

#================================ General =====================================

name: production-tools
tags: ["system"]
fields:
  sys_name: engineering-sandboxes
  sys_env: sandbox
  sys_region: us-west-2

#================================ Outputs =====================================
output.logstash:
  hosts: ["us-west-2.logstash.service.consul:5044"]
EOF
service filebeat start