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
    address = "consul.service.production.consul:8500"
    path    = "terraform/cla/application"
  }
}

# We take the State of production-tools to grab some data form there for VPC Peering Connection
data "terraform_remote_state" "cla-env" {
  backend = "consul"
  config {
    address = "consul.service.production.consul:8500"
    path    = "terraform/cla/environment"
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
data "template_file" "user_data_cla" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    env               = "${terraform.env == "default" ? "production" : terraform.env}"
    build             = "${var.build_hash}"
    ecs_cluster_name  = "${terraform.env == "default" ? "production" : terraform.env}-cla"
    region            = "${data.terraform_remote_state.cla-env.region}"
    newrelic_key      = "${data.terraform_remote_state.cla-env.newrelic_key}"
  }
}

# ECS Cluster
module "tools-ecs-cluster" {
  source                 = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/ecs-cluster"
  environment            = "${terraform.env == "default" ? "production" : terraform.env}-cla"
  team                   = "Engineering"
  name                   = "${terraform.env == "default" ? "production" : terraform.env}-cla"
  vpc_id                 = "${data.terraform_remote_state.cla-env.vpc_id}"
  subnet_ids             = "${data.terraform_remote_state.cla-env.internal_subnets}"
  key_name               = "production-cla"
  iam_instance_profile   = "${data.terraform_remote_state.cla-env.iam_profile_ecsInstance}"
  region                 = "${data.terraform_remote_state.cla-env.region}"
  availability_zones     = "${data.terraform_remote_state.cla-env.availability_zones}"
  instance_type          = "t2.medium"
  security_group         = "${data.terraform_remote_state.cla-env.sg_ecs_cluster}"
  instance_ebs_optimized = false
  desired_capacity       = "3"
  min_size               = "3"
  cloud_config_content   = "${data.template_file.user_data_cla.rendered}"
}

# Registrator
module "registrator" {
  source           = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/prod-registrator"

  # Application Information
  build_hash       = "${var.build_hash}"
  project          = "cla"

  region           = "${data.terraform_remote_state.cla-env.region}"
  ecs_cluster_name = "${module.tools-ecs-cluster.name}"
  dns_servers      = "${data.terraform_remote_state.cla-env.dns_servers}"
}

# Consul Agent
module "consul" {
  source           = "git::ssh://git@github.linuxfoundation.org/Engineering/terraform.git//modules/prod-consul-agent"

  # Consul
  encryption_key   = "9F2n4KWdxSj2Z4MMVqbHqg=="
  datacenter       = "production"
  endpoint         = "consul.service.consul"

  # Application Information
  build_hash       = "${var.build_hash}"
  project          = "cla"

  region           = "${data.terraform_remote_state.cla-env.region}"
  ecs_cluster_name = "${module.tools-ecs-cluster.name}"
  dns_servers      = "${data.terraform_remote_state.cla-env.dns_servers}"
}

# cla
module "cla" {
  source            = "cla"

  # Application Information
  build_hash      = "${var.build_hash}"
  route53_zone_id   = "${data.terraform_remote_state.cla-env.route53_zone_id}"

  # ECS Information
  internal_elb_sg   = "${data.terraform_remote_state.cla-env.sg_internal_elb}"
  internal_subnets  = "${data.terraform_remote_state.cla-env.internal_subnets}"
  region            = "${data.terraform_remote_state.cla-env.region}"
  vpc_id            = "${data.terraform_remote_state.cla-env.vpc_id}"
  ecs_cluster_name  = "${module.tools-ecs-cluster.name}"
  dns_servers       = "${data.terraform_remote_state.cla-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.cla-env.iam_role_ecsService}"
}

# NGINX Proxy
module "nginx" {
  source            = "./nginx"

  # Application Information
  build_hash      = "${var.build_hash}"
  route53_zone_id   = "${data.terraform_remote_state.cla-env.route53_zone_id}"

  # ECS Information
  external_elb_sg   = "${data.terraform_remote_state.cla-env.sg_external_elb}"
  external_subnets  = "${data.terraform_remote_state.cla-env.external_subnets}"
  region            = "${data.terraform_remote_state.cla-env.region}"
  vpc_id            = "${data.terraform_remote_state.cla-env.vpc_id}"
  ecs_cluster_name  = "${module.tools-ecs-cluster.name}"
  dns_servers       = "${data.terraform_remote_state.cla-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.cla-env.iam_role_ecsService}"
}