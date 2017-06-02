variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "build_hash" {
  description = "The Build Number we are to deploy."
}

# We are saving the state for this infra in Consul
terraform {
  backend "consul" {
    address = "consul.service.consul:8500"
    path    = "terraform/applications/pmc/app"
  }
}

# We take the State of production-tools to grab some data form there for VPC Peering Connection
data "terraform_remote_state" "pmc-env" {
  backend = "consul"
  config {
    address = "consul.service.consul:8500"
    path    = "terraform/applications/pmc/environment"
  }
}

# Provider for this infra
provider "aws" {
  alias = "local"
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

# User Data for the ECS Container Instances
data "template_file" "user_data_pmc" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    env               = "${terraform.env}"
    build             = "${var.build_hash}"
    ecs_cluster_name  = "${terraform.env}-pmc"
    region            = "${data.terraform_remote_state.pmc-env.region}"
    newrelic_key      = "${data.terraform_remote_state.pmc-env.newrelic_key}"
  }
}

# ECS Cluster
module "pmc-ecs-cluster" {
  source                 = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/ecs-cluster"
  environment            = "${terraform.env}-pmc"
  team                   = "Engineering"
  name                   = "${terraform.env}-pmc"
  vpc_id                 = "${data.terraform_remote_state.pmc-env.vpc_id}"
  subnet_ids             = "${data.terraform_remote_state.pmc-env.internal_subnets}"
  key_name               = "production-pmc"
  iam_instance_profile   = "${data.terraform_remote_state.pmc-env.iam_profile_ecsInstance}"
  region                 = "${data.terraform_remote_state.pmc-env.region}"
  availability_zones     = "${data.terraform_remote_state.pmc-env.availability_zones}"
  instance_type          = "t2.medium"
  security_group         = "${data.terraform_remote_state.pmc-env.sg_ecs_cluster}"
  instance_ebs_optimized = false
  desired_capacity       = "3"
  min_size               = "3"
  cloud_config_content   = "${data.template_file.user_data_pmc.rendered}"
}

# Registrator
module "registrator" {
  source           = "./registrator"

  # Application Information
  build_hash      = "${var.build_hash}"

  region           = "${data.terraform_remote_state.pmc-env.region}"
  ecs_cluster_name = "${module.pmc-ecs-cluster.name}"
  dns_servers      = "${data.terraform_remote_state.pmc-env.dns_servers}"
}

# Consul Agent
module "consul" {
  source           = "./consul-agent"

  # Consul
  encryption_key   = "9F2n4KWdxSj2Z4MMVqbHqg=="
  datacenter       = "AWS"

  # Application Information
  build_hash     = "${var.build_hash}"

  region           = "${data.terraform_remote_state.pmc-env.region}"
  ecs_cluster_name = "${module.pmc-ecs-cluster.name}"
  dns_servers      = "${data.terraform_remote_state.pmc-env.dns_servers}"
}

# CINCO
module "pmc" {
  source            = "./pmc"

  # Application Information
  build_hash      = "${var.build_hash}"
  route53_zone_id   = "${data.terraform_remote_state.pmc-env.route53_zone_id}"

  # ECS Information
  internal_elb_sg   = "${data.terraform_remote_state.pmc-env.sg_internal_elb}"
  internal_subnets  = "${data.terraform_remote_state.pmc-env.internal_subnets}"
  region            = "${data.terraform_remote_state.pmc-env.region}"
  vpc_id            = "${data.terraform_remote_state.pmc-env.vpc_id}"
  ecs_cluster_name  = "${module.pmc-ecs-cluster.name}"
  dns_servers       = "${data.terraform_remote_state.pmc-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.pmc-env.iam_role_ecsService}"
}

# NGINX Proxy
module "nginx" {
  source            = "./nginx"

  # Application Information
  build_hash      = "${var.build_hash}"
  route53_zone_id   = "${data.terraform_remote_state.pmc-env.route53_zone_id}"

  # ECS Information
  external_elb_sg   = "${data.terraform_remote_state.pmc-env.sg_external_elb}"
  external_subnets  = "${data.terraform_remote_state.pmc-env.external_subnets}"
  region            = "${data.terraform_remote_state.pmc-env.region}"
  vpc_id            = "${data.terraform_remote_state.pmc-env.vpc_id}"
  ecs_cluster_name  = "${module.pmc-ecs-cluster.name}"
  dns_servers       = "${data.terraform_remote_state.pmc-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.pmc-env.iam_role_ecsService}"
}