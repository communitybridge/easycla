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

data "terraform_remote_state" "cla-env" {
  backend = "consul"
  config {
    address = "consul.service.production.consul:8500"
    path    = "terraform/cla/environment"
  }
}

data "terraform_remote_state" "cla-app" {
  backend = "consul"
  config {
    address = "consul.service.production.consul:8500"
    path    = "terraform/cla/application"
  }
}

# Provider for this infra
provider "aws" {
  alias = "local"
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}


# CLA console
module "cla-console" {
  source            = "console"

  # Application Information
  build_hash        = "${var.build_hash}"
  route53_zone_id   = "${data.terraform_remote_state.cla-env.route53_zone_id}"

  # ECS Information
  internal_elb_sg   = "${data.terraform_remote_state.cla-env.sg_internal_elb}"
  internal_subnets  = "${data.terraform_remote_state.cla-env.internal_subnets}"
  region            = "${data.terraform_remote_state.cla-env.region}"
  vpc_id            = "${data.terraform_remote_state.cla-env.vpc_id}"
  ecs_cluster_name  = "production-cla"
  dns_servers       = "${data.terraform_remote_state.cla-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.cla-env.iam_role_ecsService}"
}

# NGINX Proxy
module "nginx" {
  source            = "./nginx"

  # Application Information
  build_hash        = "${var.build_hash}"
  route53_zone_id   = "cla_route53"

  # ECS Information
  external_elb_sg   = "${data.terraform_remote_state.cla-env.sg_external_elb}"
  external_subnets  = "${data.terraform_remote_state.cla-env.external_subnets}"
  region            = "${data.terraform_remote_state.cla-env.region}"
  vpc_id            = "${data.terraform_remote_state.cla-env.vpc_id}"
  ecs_cluster_name  = "production-cla"
  dns_servers       = "${data.terraform_remote_state.cla-env.dns_servers}"
  ecs_role          = "${data.terraform_remote_state.cla-env.iam_role_ecsService}"
}