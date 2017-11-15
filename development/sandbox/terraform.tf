variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "cidr" {
  default = "10.32.1.0/24"
}

provider "aws" {
  region     = "us-west-2"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

terraform {
    backend "consul" {
      address = "consul.service.development.consul:8500"
      path    = "terraform/sandboxes"
    }
}

// This allows me to pull the state of another environment, in this case production-tools and grab data from it.
data "terraform_remote_state" "infrastructure" {
  backend = "consul"
  config {
    address = "consul.service.development.consul:8500"
    path    = "terraform/infrastructure"
  }
}

module "peering" {
  source                    = "../../modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.infrastructure.account_number}"
  tools_cidr                = "${data.terraform_remote_state.infrastructure.west_cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.infrastructure.west_vpc_id}"
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "sandboxes"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "us-west-2"
  }
}

module "vpc" {
  source             = "../../modules/vpc"
  name               = "Sandboxes"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.32.1.128/27", "10.32.1.160/27", "10.32.1.192/27"]
  external_subnets   = ["10.32.1.0/27",   "10.32.1.32/27",  "10.32.1.64/27"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]
}

module "dhcp" {
  source  = "../../modules/dhcp"
  name    = "ci.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "10.32.0.140, 10.32.0.180, 10.32.0.220"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "sandboxes"
}

module "engineering-sandboxes-ecs-cluster" {
  source               = "../../modules/ecs-cluster"
  environment          = "Sandbox"
  team                 = "Engineering"
  name                 = "sandboxes"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  key_name             = "engineering-sandboxes"
  iam_instance_profile = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region               = "us-west-2"
  availability_zones   = "${module.vpc.availability_zones}"
  instance_type        = "t2.medium"
  security_group       = "${module.security_groups.engineering_sandboxes}"
  instance_ebs_optimized = false
  desired_capacity     = "3"
  min_size             = "3"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
}

module "rds-cluster" {
  source               = "../../modules/rds-cluster"
  master_username      = "lfengineering"
  name                 = "sandboxes"
  master_password      = "buanCAWwwAGxUyoU2Fai"
  availability_zones   = "${module.vpc.availability_zones}"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  environment          = "Sandbox"
  team                 = "Engineering"
  security_groups      = ["${module.engineering-sandboxes-ecs-cluster.security_group_id}"]
  engine               = "mariadb"
  engine_version       = "10.1.19"
  parameter_group_name = "engineering"
  instance_type        = "db.t2.medium"
}

module "redis-cluster" {
  source               = "../../modules/redis-cluster"
  name                 = "sandboxes"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  environment          = "Sandbox"
  team                 = "Engineering"
  security_groups      = ["${module.security_groups.engineering_sandboxes_redis}"]
  instance_type        = "cache.t2.small"
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
  value = "${module.security_groups.engineering_sandboxes_elb}"
}

output "cidr" {
  value = "${var.cidr}"
}

output "vpc_id" {
  value = "${module.vpc.id}"
}

output "ecs_cluster_name" {
  value = "${module.engineering-sandboxes-ecs-cluster.name}"
}

output "internal_elb_sg" {
  value = "${module.security_groups.engineering_sandboxes_elb}"
}

output "external_elb_sg" {
  value = "${module.security_groups.engineering_sandboxes_elb}"
}

output "sandbox_cert_arn" {
  value = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
}

output "ecs_role" {
  value = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
}
