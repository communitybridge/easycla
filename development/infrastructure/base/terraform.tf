variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "name" {
  description = "Name of your VPC"
}

variable "region" {
  description = "Region to launch this infra on"
}

variable "cidr" {
  description = "The CIDR block for the VPC."
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = "list"
}

variable "external_subnets" {
  description = "List of external subnets"
  type        = "list"
}

variable "internal_subnets" {
  description = "List of internal subnets"
  type        = "list"
}

variable "key_name" {
  description = "Key Pair to use to administer this vpc"
}

variable "newrelic_key" {
  description = "Key to use for NewRelic"
}

variable "region_identitier" {
  description = "Label to recognize the region"
}

variable "consul_encryption_key" {
  description = "The encryption key used for consul"
}

variable "dns_server" {
  description = "DNS Server for the VPC"
}

variable "r53_zone_id" {}

variable "ghe_peering" {}

variable "nexus" {}

provider "aws" {
  region     = "${var.region}"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "aws_cloudwatch_log_group" "infrastructure" {
  name = "infrastructure"
}

module "vpc" {
  source             = "../../modules/vpc"
  name               = "${var.name}"
  cidr               = "${var.cidr}"
  internal_subnets   = "${var.internal_subnets}"
  external_subnets   = "${var.external_subnets}"
  availability_zones = "${var.availability_zones}"
}

module "dhcp" {
  source  = "../../modules/dhcp"
  name    = "prod.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "${replace(var.cidr, ".0/24", ".140")},${replace(var.cidr, ".0/24", ".180")},${replace(var.cidr, ".0/24", ".220")},${cidrhost(var.cidr, 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "Engineering"
}

resource "aws_route53_zone_association" "prod_zone" {
  zone_id      = "${var.r53_zone_id}"
  vpc_id       = "${module.vpc.id}"
  vpc_region   = "${var.region}"
}

resource "aws_vpc_peering_connection" "peer" {
  provider      = "aws.local"
  count         = "${var.ghe_peering}"

  peer_owner_id = "961082193871"
  peer_vpc_id   = "vpc-10c9f477"
  vpc_id        = "${module.vpc.id}"

  accepter {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "peer_internal_1" {
  provider                  = "aws.local"
  count                     = "${var.ghe_peering}"
  route_table_id            = "${module.vpc.raw_route_tables_id[0]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_2" {
  provider                  = "aws.local"
  count                     = "${var.ghe_peering}"
  route_table_id            = "${module.vpc.raw_route_tables_id[1]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_3" {
  provider                  = "aws.local"
  count                     = "${var.ghe_peering}"
  route_table_id            = "${module.vpc.raw_route_tables_id[2]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.local"
  count                     = "${var.ghe_peering}"
  route_table_id            = "${module.vpc.external_rtb_id}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

# Creating EFS for Tools Storage
resource "aws_efs_file_system" "production-tools-storage" {
  provider = "aws.local"
  creation_token = "production-tools-storage"

  tags {
    Name = "Enginnering - Production Tools Storage"
  }
}

resource "aws_efs_mount_target" "efs_mount_1" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.production-tools-storage.id}"
  subnet_id       = "${module.vpc.internal_subnets[0]}"
  security_groups = ["${module.security_groups.efs}"]
}

resource "aws_efs_mount_target" "efs_mount_2" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.production-tools-storage.id}"
  subnet_id       = "${module.vpc.internal_subnets[1]}"
  security_groups = ["${module.security_groups.efs}"]
}

resource "aws_efs_mount_target" "efs_mount_3" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.production-tools-storage.id}"
  subnet_id       = "${module.vpc.internal_subnets[2]}"
  security_groups = ["${module.security_groups.efs}"]
}

data "template_file" "ecs_instance_cloudinit_tools" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name  = "infrastructure"
    region            = "${var.region}"
    region_identifier = "${var.region_identitier}"
    newrelic_key      = "${var.newrelic_key}"
    efs_id            = "${aws_efs_file_system.production-tools-storage.id}"
  }
}

module "tools-ecs-cluster" {
  source                 = "../../modules/ecs-cluster"
  environment            = "Production"
  team                   = "Engineering"
  name                   = "infrastructure"
  vpc_id                 = "${module.vpc.id}"
  subnet_ids             = "${module.vpc.internal_subnets}"
  key_name               = "${var.key_name}"
  iam_instance_profile   = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region                 = "${var.region}"
  availability_zones     = "${module.vpc.availability_zones}"
  instance_type          = "t2.large"
  security_group         = "${module.security_groups.tools-ecs-cluster}"
  instance_ebs_optimized = false
  desired_capacity       = "1"
  min_size               = "1"
  cloud_config_content   = "${data.template_file.ecs_instance_cloudinit_tools.rendered}"
}

module "consul-bind" {
  source                 = "./consul_bind"

  internal_subnets       = "${module.vpc.internal_subnets}"
  bind_sg                = "${module.security_groups.bind}"
  cidr                   = "${var.cidr}"
  region                 = "${var.region}"
  internal_elb_sg        = "${module.security_groups.internal_elb}"
}

module "registrator" {
  source = "./registrator"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  dns_servers            = "${module.consul-bind.dns_servers}"
  region                 = "${var.region}"
}

module "consul-agent" {
  source = "./consul-agent"

  encryption_key   = "${var.consul_encryption_key}"
  datacenter       = "development"
  endpoint         = "consul.service.development.consul"
  ecs_cluster_name = "${module.tools-ecs-cluster.name}"
  dns_servers      = "${module.consul-bind.dns_servers}"
}

module "nexus" "nexus-master" {
  source                 = "./nexus"

  building               = "${var.nexus}"
  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  internal_subnets       = "${module.vpc.internal_subnets}"
  internal_elb_sg        = "${module.security_groups.internal_elb}"
  dns_servers            = "${module.consul-bind.dns_servers}"

  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"
}

module "pritunl" {
  source                 = "./pritunl"

  external_subnets       = "${module.vpc.external_subnets}"
  vpn_sg                 = "${module.security_groups.vpn}"
}

module "mongodb" {
  source = "./mongodb"

  vpc_id = "${module.vpc.id}"
  region = "${var.region}"
  dns_servers = "${module.consul-bind.dns_servers}"
  ecs_cluster_name = "${module.tools-ecs-cluster.name}"
}

module "it-managed-vpn" {
  source            = "./it_vpn_tunnel"

  vpc_id            = "${module.vpc.id}"
  internal_subnets  = "${module.vpc.internal_subnets}"
  cidr              = "10.32.0.0/12"
  key_name          = "${var.key_name}"
  route_tables      = "${module.vpc.raw_route_tables_id}"
}

module "cert-managed-vpn" {
  source            = "./cert_vpn_tunnel"

  vpc_id            = "${module.vpc.id}"
  external_subnets  = "${module.vpc.external_subnets}"
  cidr              = "10.32.0.0/12"
  key_name          = "${var.key_name}"
  route_tables      = "${module.vpc.raw_route_tables_id}"
}

# RDS Instance to use by some services, ie: Keycloak
module "rds-cluster" {
  source               = "../../modules/rds-cluster"
  master_username      = "lfengineering"
  name                 = "keycloak"
  master_password      = "buanCAWwwAGxUyoU2Fai"
  availability_zones   = "${var.availability_zones}"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  environment          = "Preprod"
  team                 = "Engineering"
  security_groups      = ["${module.security_groups.tools-ecs-cluster}"]
  engine               = "mariadb"
  engine_version       = "10.1.19"
  parameter_group_name = "engineering"
  instance_type        = "db.t2.small"
}

# LDAP used for Keycloak sandbox
module "open-ldap" {
  source                 = "./openldap"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  internal_subnets       = "${var.internal_subnets}"
  internal_elb_sg        = "${module.tools-ecs-cluster.security_group_id}"
  dns_servers            = ["10.32.0.140", "10.32.0.180", "10.32.0.220"]

  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"

  ldap_org               = "linuxfoundation"
  ldap_domain            = "linuxfoundation.org"
  ldap_admin_password    = "ZPw4RRzxLikVdN"
}

module "keycloak" {
  source                 = "./keycloak"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  subnets                = "${var.external_subnets}"
  internal_elb_sg        = "${module.security_groups.internal_elb}"
  dns_servers            = ["10.32.0.140", "10.32.0.180", "10.32.0.220"]

  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"
  env                    = "sandbox"

  mysql_db               = "sandbox"
  mysql_host             = "keycloak.cnfn2tun3mjw.us-west-2.rds.amazonaws.com"
  mysql_pass             = "buanCAWwwAGxUyoU2Fai"
  mysql_port             = "3306"
  mysql_user             = "lfengineering"
}


/**
 * Outputs
 */


// The region in which the infra lives.
output "region" {
  value = "${var.region}"
}

// The VPC's CIDR
output "cidr" {
  value = "${var.cidr}"
}

// Comma separated list of internal subnet IDs.
output "internal_subnets" {
  value = "${module.vpc.internal_subnets}"
}

// Comma separated list of external subnet IDs.
output "external_subnets" {
  value = "${module.vpc.external_subnets}"
}

// The VPC availability zones.
output "availability_zones" {
  value = "${module.vpc.availability_zones}"
}

// The VPC security group ID.
output "vpc_security_group" {
  value = "${module.vpc.security_group}"
}

// The VPC ID.
output "vpc_id" {
  value = "${module.vpc.id}"
}

// The VPC security group ID.
output "raw_route_tables_id" {
  value = "${module.vpc.raw_route_tables_id}"
}

// The VPC ID.
output "external_rtb_id" {
  value = "${module.vpc.external_rtb_id}"
}

// Comma separated list of internal route table IDs.
output "internal_route_tables" {
  value = "${module.vpc.internal_rtb_id}"
}

// The external route table ID.
output "external_route_tables" {
  value = "${module.vpc.external_rtb_id}"
}

// External ELB allows traffic from the world.
output "sg_internal_elb" {
  value = "${module.security_groups.internal_elb}"
}

output "tools_ecs_asg" {
  value = "${module.tools-ecs-cluster.asg_name}"
}

output "tools_ecs_name" {
  value = "${module.tools-ecs-cluster.name}"
}

output "tools_ecs_sg" {
  value = "${module.tools-ecs-cluster.security_group_id}"
}

//output "consul_elb_cname" {
//  value = "${module.consul.consul_elb_cname}"
//}
//
//output "consul_elb_name" {
//  value = "${module.consul.consul_elb_name}"
//}
//
//output "consul_service_name" {
//  value = "${module.consul.consul_service_name}"
//}
//
//output "consul_elb_zoneid" {
//  value = "${module.consul.consul_elb_zoneid}"
//}

output "dns_servers" {
  value = "${module.consul-bind.dns_servers}"
}
