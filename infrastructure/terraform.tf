variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "region" {
  default = "us-west-2"
}

variable "external_dns_name" {
  default = "e.tux.rocks"
}

variable "external_dns_zoneid" {
  default = "Z1CCIH5LOWNGEU"
}

provider "aws" {
  region     = "${var.region}"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

module "vpc" {
  source             = "../modules/vpc"
  name               = "Infrastructure"
  cidr               = "10.45.0.0/16"
  internal_subnets   = ["10.45.0.0/19", "10.45.32.0/19", "10.45.64.0/19"]
  external_subnets   = ["10.45.96.0/19", "10.45.128.0/19", "10.45.160.0/19"]
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "aws.eng.tux.rocks"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost("10.45.0.0/16", 2)}"
}

module "dns" {
  source = "../modules/dns"
  name   = "aws.eng.tux.rocks"
  vpc_id = "${module.vpc.id}"
}

module "security_groups" "sgs" {
  source = "security_groups/"
  vpc_id = "${module.vpc.id}"
}

module "s3_buckets" {
  source = "s3_buckets"
}

module "pritunl" {
  source            = "modules/pritunl"

  external_subnets  = "${module.vpc.external_subnets}"
  internal_subnets  = "${module.vpc.internal_subnets}"
  sg                = "${module.security_groups.pritunl_node}"
  elb_sg            = "${module.security_groups.pritunl_elb}"
  mongodb_sg        = "${module.security_groups.pritunl_mongodb}"
  dns_name          = "${var.external_dns_name}"
  dns_zone_id       = "${var.external_dns_zoneid}"
}
