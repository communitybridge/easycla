variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

provider "aws" {
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

module "vpc" "engineering_vpc" {
  source             = "../modules/vpc"
  name               = "Engineering"
  cidr               = "10.40.0.0/16"
  internal_subnets   = ["10.40.0.0/19" ,"10.40.64.0/19", "10.40.128.0/19"]
  external_subnets   = ["10.40.32.0/20", "10.40.96.0/20", "10.40.160.0/20"]
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

module "dns" {
  source = "../modules/dns"
  name   = "engineering.local"
  vpc_id = "${module.vpc.id}"
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "${module.dns.name}"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost("10.40.0.0/16", 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "10.40.0.0/16"
  vpc_id  = "${module.vpc.id}"
  name    = "engineering"
}

module "bastion" {
  source          = "../modules/bastion"
  region          = "us-west-2"
  instance_type   = "t2.micro"
  security_groups = "${module.security_groups.external_ssh},${module.security_groups.internal_ssh}"
  vpc_id          = "${module.vpc.id}"
  subnet_id       = "${element(module.vpc.external_subnets, 0)}"
  key_name        = "engineering-tools"
  instance_name   = "Bastion: Engineering"
}

// The region in which the infra lives.
output "region" {
  value = "us-west-2"
}

// The internal route53 zone ID.
output "zone_id" {
  value = "${module.dns.zone_id}"
}

// The VPC's CIDR
output "cidr" {
  value = "10.40.0.0/16"
}

// Comma separated list of internal subnet IDs.
output "internal_subnets" {
  value = "${module.vpc.internal_subnets}"
}

// Comma separated list of external subnet IDs.
output "external_subnets" {
  value = "${module.vpc.external_subnets}"
}

// The internal domain name, e.g "stack.local".
output "domain_name" {
  value = "${module.dns.name}"
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

// Comma separated list of internal route table IDs.
output "internal_route_tables" {
  value = "${module.vpc.internal_rtb_id}"
}

// The external route table ID.
output "external_route_tables" {
  value = "${module.vpc.external_rtb_id}"
}

// External SSH allows ssh connections on port 22 from the world.
output "sg_external_ssh" {
  value = "${module.security_groups.external_ssh}"
}

// Internal SSH allows ssh connections from the external ssh security group.
output "sg_internal_ssh" {
  value = "${module.security_groups.internal_ssh}"
}

// External ELB allows traffic from the world.
output "sg_external_elb" {
  value = "${module.security_groups.external_elb}"
}

// External ELB allows traffic from the world.
output "sg_internal_elb" {
  value = "${module.security_groups.internal_elb}"
}

// External ELB allows traffic from the world.
output "sg_vpn" {
  value = "${module.security_groups.vpn}"
}

