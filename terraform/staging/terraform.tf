variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "cidr" {
  default = "10.37.0.0/16"
}

provider "aws" {
  region     = "us-west-2"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

terraform {
  backend "consul" {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/staging"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

// This allows me to pull the state of another environment, in this case production-tools and grab data from it.
data "terraform_remote_state" "infrastructure" {
  backend = "consul"
  config {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/infrastructure2.0"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

resource "aws_cloudwatch_log_group" "staging" {
  provider = "aws.local"
  name = "staging"
}

module "peering" {
  source                    = "../modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.infrastructure.account_number}"
  tools_cidr                = "${data.terraform_remote_state.infrastructure.cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.infrastructure.vpc_id}"
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "staging"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "us-west-2"
  }
}

module "vpc" {
  source             = "../modules/vpc"
  name               = "Staging Applications"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.37.0.0/22", "10.37.4.0/22", "10.37.8.0/22"]
  external_subnets   = ["10.37.12.0/22",   "10.37.16.0/22",  "10.37.20.0/22"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "staging.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost(var.cidr, 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "sandboxes"
}

// IAM Profiles, Roles & Policies for ECS
module "ecs-iam-profile" {
  source = "./iam-role"

  name = "cluster"
  environment = "staging"
}

module "ecs-cluster" {
  source               = "../modules/ecs-cluster"
  environment          = "Staging"
  team                 = "Engineering"
  name                 = "staging"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  key_name             = "engineering-staging"
  iam_instance_profile = "${module.ecs-iam-profile.ecsInstanceProfile}"
  region               = "us-west-2"
  availability_zones   = "${module.vpc.availability_zones}"
  instance_type        = "m5.large"
  security_group       = "${module.security_groups.engineering_staging}"
  instance_ebs_optimized = false
  desired_capacity     = "4"
  min_size             = "4"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
  root_volume_size     = 100
  docker_volume_size   = 100
}

# Consul Agent
module "consul" {
  source           = "../modules/consul-agent"

  # Consul
  encryption_key   = "yYQCMdOCtX73dmYeEQ/NYA=="
  datacenter       = "us-west-2"
  endpoint         = "consul.eng.linuxfoundation.org"
  ecs_cluster_name = "staging"
  log_group_name   = "staging"
}

# Consul Agent
module "registrator" {
  source           = "../modules/registrator"

  # Consul
  ecs_cluster_name   = "staging"
  project            = "staging"
  log_group_name     = "staging"
}

output "internal_subnets" {
  value = "${module.vpc.internal_subnets}"
}

output "external_subnets" {
  value = "${module.vpc.external_subnets}"
}

output "raw_route_tables_id" {
  value = "${module.vpc.raw_route_tables_id}"
}

output "sg_external_elb" {
  value = "${module.security_groups.engineering_staging_elb}"
}

output "cidr" {
  value = "${var.cidr}"
}

output "vpc_id" {
  value = "${module.vpc.id}"
}

output "ecs_cluster_name" {
  value = "${module.ecs-cluster.name}"
}

output "ecs_cluster_security_group" {
  value = "${module.ecs-cluster.security_group_id}"
}