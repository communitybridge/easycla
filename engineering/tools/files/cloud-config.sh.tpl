#!/bin/bash

# Join the default ECS cluster
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

PATH=$PATH:/usr/local/bin

# Instance should be added to an security group that allows HTTP outbound
yum update -y

# Install packages
yum -y install jq nfs-utils python27 python27-pip

# Get region of EC2 from instance metadata
EC2_AVAIL_ZONE=`curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone`
EC2_REGION="`echo \"$EC2_AVAIL_ZONE\" | sed -e 's:\([0-9][0-9]*\)[a-z]*\$:\\1:'`"

# Create mount point
mkdir /mnt/storage

# Get EFS FileSystemID attribute
#Instance needs to be added to a EC2 role that give the instance at least read access to EFS
EFS_FILE_SYSTEM_ID=${efs_id}

# Instance needs to be a member of security group that allows 2049 inbound/outbound
#The security group that the instance belongs to has to be added to EFS file system configuration
#Create variables for source and target
DIR_SRC=$EC2_AVAIL_ZONE.$EFS_FILE_SYSTEM_ID.efs.$EC2_REGION.amazonaws.com:/
DIR_TGT=/mnt/storage

# Mount EFS file system
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $DIR_SRC $DIR_TGT

#Backup fstab
cp -p /etc/fstab /etc/fstab.back-$(date +%F)

#Append line to fstab
echo -e "$DIR_SRC \t\t $DIR_TGT \t\t nfs \t\t defaults \t\t 0 \t\t 0" | tee -a /etc/fstab

EC2_INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`

#NewRelic Installation
rpm -Uvh https://download.newrelic.com/pub/newrelic/el5/i386/newrelic-repo-5-3.noarch.rpm
yum install newrelic-sysmond -y
usermod -a -G docker newrelic
nrsysmond-config --set license_key=${newrelic_key}
sed -i '/#hostname*/c\hostname=\"${newrelic_hostname} ('$EC2_INSTANCE_ID')\"' /etc/newrelic/nrsysmond.cfg
sed -i '/#labels*/c\labels=${newrelic_labels}' /etc/newrelic/nrsysmond.cfg
/etc/init.d/newrelic-sysmond start

# Setting for ES
sysctl -w vm.max_map_count=262144

# Using a different CIDR for Docker (not to conflict with coreIT VPC CIDR)
sed -i '$ d' /etc/sysconfig/docker
echo 'OPTIONS="--default-ulimit nofile=1024:4096 --bip=192.168.5.1/24"' >> /etc/sysconfig/docker
service docker stop
ip link set docker0 down
ip link delete docker0
service docker start

# Restarting ECS
start ecs
