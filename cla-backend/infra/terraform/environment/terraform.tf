variable "access_key" {
  description = "The Terraform Production AWS Access Key (LastPass)"
}

variable "secret_key" {
  description = "The Terraform Production Your AWS Secret Key (LastPass)"
}

variable "cidr" {
  default     = "10.32.5.0/24"
  description = "The CIDR used to bootstrap this new environment"
}

// We are using consul as our storage backend for the terraform state, it supports locking natively.
terraform {
  backend "consul" {
    address = "consul.service.production.consul:8500"
    path    = "terraform/cla/environment"
  }
}

// This allows me to pull the state of another environment, in this case production-tools and grab data from it.
data "terraform_remote_state" "production-tools" {
  backend = "consul"
  config {
    address = "consul.service.development.consul:8500"
    path    = "terraform/infrastructure"
  }
}

// This allows me to pull the state of another environment, in this case production-tools and grab data from it.
data "terraform_remote_state" "infrastructure" {
  backend = "consul"
  config {
    address = "consul.service.production.consul:8500"
    path    = "terraform/infrastructure"
  }
}

// Provider for this infra, re-using the same credentials that asked for in the variables.
provider "aws" {
  alias = "local"
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

// Provisioning the VPC for the environment, please do not edit the CIDR allocations without looking at Confluence
// inside the CIDR Allocation section. Those Subnets have been pre-approved in other infras. Don't touch! :)
module "vpc" {
  source             = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/vpc"
  name               = "CLA"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.32.5.128/27", "10.32.5.160/27", "10.32.5.192/27"]
  external_subnets   = ["10.32.5.0/27",   "10.32.5.32/27",  "10.32.5.64/27"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]
}

// Making sure that each Internal Subnet has a NAT attached to it, with the proper DNS Servers set in the
// DHCP options inside the VPC. We are using Production Tools's DNS Servers to resolve to the internet.
module "dhcp" {
  source  = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/dhcp"
  name    = "cla.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "${join(",", data.terraform_remote_state.infrastructure.west_dns_servers)}"
}

// Holds all the security groups for this infra and the application. Edit with caution.
module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "cla"
}

// App: Repository to hold cla production artifacts.
resource "aws_ecr_repository" "production" {
  provider = "aws.local"
  name = "cla/production"
}

// App: Repository to hold NGINX production artifacts.
resource "aws_ecr_repository" "nginx" {
  provider = "aws.local"
  name = "cla/nginx"
}

// App: Repository to hold CLA Console production artifacts.
resource "aws_ecr_repository" "console_production" {
  provider = "aws.local"
  name = "cla-console/production"
}

// App: Repository to hold NGINX CLA Console production artifacts.
resource "aws_ecr_repository" "console_nginx" {
  provider = "aws.local"
  name = "cla-console/nginx"
}

resource "aws_s3_bucket" "assets_bucket" {
  provider  = "aws.local"
  bucket    = "lf-cla-artifacts"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "GET", "DELETE"]
    allowed_origins = ["*"]
    max_age_seconds = 3000
  }

  versioning {
    enabled = true
  }
}

// IAM Profiles, Roles & Policies for ECS
module "ecs-iam-profile" {
  source = "./iam-role"

  name = "cla"
  environment = "production"
}

// Peering to Production Tools
module "peering" {
  source                    = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.production-tools.account_number}"
  tools_cidr                = "${data.terraform_remote_state.production-tools.west_cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.production-tools.west_vpc_id}"
}

// Peering to Infrastructure
module "infra_peering" {
  source                    = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.infrastructure.account_number}"
  tools_cidr                = "${data.terraform_remote_state.infrastructure.west_cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.infrastructure.west_vpc_id}"
}

// The region in which the infra lives.
output "region" {
  value = "us-west-2"
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

// SG: ECS-Cluster
output "sg_ecs_cluster" {
  value = "${module.security_groups.ecs-cluster}"
}

// SG: External ELB
output "sg_external_elb" {
  value = "${module.security_groups.external_elb}"
}

// SG: Internal ELB
output "sg_internal_elb" {
  value = "${module.security_groups.external_elb}"
}

// ecsService role to re-use inside the App Terraform Script
output "iam_role_ecsService" {
  value = "${module.ecs-iam-profile.service_role_arn}"
}

// ecsInstanceProfile role to re-use inside the App Terraform Script
output "iam_profile_ecsInstance" {
  value = "${module.ecs-iam-profile.ecsInstanceProfile}"
}

// NewRelic License Key
output "newrelic_key" {
  value = "${data.terraform_remote_state.production-tools.newrelic_key}"
}

// DNS Servers from Production-Tools
output "dns_servers" {
  value = "${data.terraform_remote_state.infrastructure.west_dns_servers}"
}

output "api_route53" {
  value = "Z2RO1I6G276B3X"
}

output "cla_route53" {
  value = "Z3I9LUXIPY6WWD"
}