#!/bin/bash

yum update -y
yum remove java-1.7.0 -y
yum install -y gcc java-1.8.0-openjdk git yum-utils openssl-devel install jq nfs-utils python27 python27-pip

echo "nameserver 10.32.2.2" > /etc/resolv.conf

echo "license_key: ${newrelic_license}" | tee -a /etc/newrelic-infra.yml
sudo curl -o /etc/yum.repos.d/newrelic-infra.repo https://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'
yum install newrelic-infra -y

wget -O /etc/yum.repos.d/jenkins.repo http://pkg.jenkins-ci.org/redhat-stable/jenkins.repo
rpm --import http://pkg.jenkins-ci.org/redhat-stable/jenkins-ci.org.key
yum install jenkins -y

DIR_SRC=/dev/nvme1n1
DIR_TGT=/var/lib/jenkins

# Mount EFS file system
mount $DIR_SRC $DIR_TGT

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

name: jenkins-master
tags: ["system"]
fields:
  sys_name: jenkins
  sys_env: master
  sys_region: ${aws_region}

#================================ Outputs =====================================
output.logstash:
  hosts: ["${aws_region}.logstash.service.consul:5044"]
EOF

service filebeat start

initctl restart newrelic-infra