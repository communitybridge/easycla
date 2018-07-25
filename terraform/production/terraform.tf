variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "cidr" {
  default = "10.38.0.0/16"
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
    path    = "terraform/production"
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
    ecs_cluster_name = "production"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "us-west-2"
  }
}

module "vpc" {
  source             = "../modules/vpc"
  name               = "Production Applications"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.38.0.0/22", "10.38.4.0/22", "10.38.8.0/22"]
  external_subnets   = ["10.38.12.0/22",   "10.38.16.0/22",  "10.38.20.0/22"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "production.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost(var.cidr, 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "production"
}

// IAM Profiles, Roles & Policies for ECS
module "ecs-iam-profile" {
  source = "./iam-role"

  name = "production-cluster"
  environment = "production"
}

module "engineering-production-ecs-cluster" {
  source               = "../modules/ecs-cluster"
  environment          = "Production"
  team                 = "Engineering"
  name                 = "production"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  key_name             = "engineering-production"
  iam_instance_profile = "${module.ecs-iam-profile.ecsInstanceProfile}"
  region               = "us-west-2"
  availability_zones   = "${module.vpc.availability_zones}"
  instance_type        = "m5.large"
  security_group       = "${module.security_groups.engineering_production}"
  instance_ebs_optimized = false
  desired_capacity     = "4"
  min_size             = "4"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
  root_volume_size     = 100
  docker_volume_size   = 100
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
  value = "${module.security_groups.engineering_production_elb}"
}

output "cidr" {
  value = "${var.cidr}"
}

output "vpc_id" {
  value = "${module.vpc.id}"
}

output "ecs_cluster_name" {
  value = "${module.engineering-production-ecs-cluster.name}"
}

output "ecs_cluster_security_group" {
  value = "${module.engineering-production-ecs-cluster.security_group_id}"
}
