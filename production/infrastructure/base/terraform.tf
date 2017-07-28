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

variable "lb_bucket" {}

provider "aws" {
  region     = "${var.region}"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

module "vpc" {
  source             = "../../../modules/vpc"
  name               = "${var.name}"
  cidr               = "${var.cidr}"
  internal_subnets   = "${var.internal_subnets}"
  external_subnets   = "${var.external_subnets}"
  availability_zones = "${var.availability_zones}"
}

module "dhcp" {
  source  = "../../../modules/dhcp"
  name    = "eng.linuxfoundation.org"
  vpc_id  = "${module.vpc.id}"
  servers = "${replace(var.cidr, ".0.0/21", ".0.10")},${replace(var.cidr, ".0.0/21", ".1.10")},${replace(var.cidr, ".0.0/21", ".2.10")},${cidrhost(var.cidr, 2)}"
//  servers = "${cidrhost(var.cidr, 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "Infrastructure"
}

module "iams" {
  source      = "./iams"

  name        = "infrastructure"
  environment = "prod"
}

# Creating combined zone for prod
resource "aws_route53_zone" "prod" {
  name = "${var.region_identitier}.prod-infra.linuxfoundation.internal."
  vpc_id = "${module.vpc.id}"
  vpc_region = "${var.region}"
}

resource "aws_cloudwatch_log_group" "infrastructure" {
  name = "infrastructure"
}

resource "aws_vpc_peering_connection" "sandbox" {
  provider      = "aws.local"

  peer_owner_id = "433610389961"
  peer_vpc_id   = "vpc-f4201f93"
  vpc_id        = "${module.vpc.id}"

  accepter {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "peer_internal_1" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[0]}"
  destination_cidr_block    = "10.32.0.0/24"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.sandbox.id}"
}

resource "aws_route" "peer_internal_2" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[1]}"
  destination_cidr_block    = "10.32.0.0/24"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.sandbox.id}"
}

resource "aws_route" "peer_internal_3" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[2]}"
  destination_cidr_block    = "10.32.0.0/24"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.sandbox.id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.external_rtb_id}"
  destination_cidr_block    = "10.32.0.0/24"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.sandbox.id}"
}

data "template_file" "ecs_instance_cloudinit_tools" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name  = "infrastructure"
    region            = "${var.region}"
    region_identifier = "${var.region_identitier}"
    newrelic_key      = "${var.newrelic_key}"
  }
}

module "tools-ecs-cluster" {
  source                 = "../../../modules/ecs-cluster"
  environment            = "Production"
  team                   = "Engineering"
  name                   = "infrastructure"
  vpc_id                 = "${module.vpc.id}"
  subnet_ids             = "${module.vpc.internal_subnets}"
  key_name               = "${var.key_name}"
  iam_instance_profile   = "${module.iams.ecsInstanceProfile}"
  region                 = "${var.region}"
  availability_zones     = "${module.vpc.availability_zones}"
  instance_type          = "t2.medium"
  security_group         = "${module.security_groups.tools-ecs-cluster}"
  instance_ebs_optimized = false
  desired_capacity       = "3"
  min_size               = "3"
  cloud_config_content   = "${data.template_file.ecs_instance_cloudinit_tools.rendered}"
}

module "consul-bind" {
  source                 = "./consul_bind"

  internal_subnets       = "${module.vpc.internal_subnets}"
  bind_sg                = "${module.security_groups.bind}"
  cidr                   = "${var.cidr}"
  region                 = "${var.region}"
  internal_elb_sg        = "${module.security_groups.internal_elb}"
  r53_zone_id            = "Z1WGT54F777KX0"
  keypair                = "${var.key_name}"
  iam_role               = "${module.iams.ecsInstanceProfile}"
  lb_logs_bucket         = "${var.lb_bucket}"
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
  datacenter       = "production"
  endpoint         = "consul.service.production.consul"
  ecs_cluster_name = "${module.tools-ecs-cluster.name}"
  dns_servers      = "${module.consul-bind.dns_servers}"
}

module "vault" "vault-master" {
  source                 = "./vault"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  internal_subnets       = "${module.vpc.internal_subnets}"
  internal_elb_sg        = "${module.security_groups.internal_elb}"
  consul_endpoint        = "127.0.0.1:8500"
  dns_servers            = "${module.consul-bind.dns_servers}"
  region                 = "${var.region}"
}

module "logstash" {
  source                 = "./logstash"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  dns_servers            = "${module.consul-bind.dns_servers}"
  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"
}

module "logstash-cloudwatch" {
  source                 = "./logstash-cloudwatch"

  ecs_cluster_name       = "${module.tools-ecs-cluster.name}"
  ecs_asg_name           = "${module.tools-ecs-cluster.asg_name}"
  dns_servers            = "${module.consul-bind.dns_servers}"
  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"
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

output "dns_servers" {
  value = "${module.consul-bind.dns_servers}"
}
