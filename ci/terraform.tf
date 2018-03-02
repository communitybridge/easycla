variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "cidr" {
  default = "10.32.2.0/24"
}

provider "aws" {
  region     = "us-west-2"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

provider "aws" {
  region     = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

terraform {
  backend "consul" {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/development/ci"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

// This allows me to pull the state of another environment, in this case production-tools and grab data from it.
data "terraform_remote_state" "prod_infrastructure" {
  backend = "consul"
  config {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/infrastructure"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

data "terraform_remote_state" "infrastructure2" {
  backend = "consul"
  config {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/infrastructure2.0"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

module "peering_infra" {
  source                    = "../modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.infrastructure2.account_number}"
  tools_cidr                = "${data.terraform_remote_state.infrastructure2.cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.infrastructure2.vpc_id}"
}

module "vpc" {
  source             = "../modules/vpc"
  name               = "CI"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.32.2.128/27", "10.32.2.160/27", "10.32.2.192/27"]
  external_subnets   = ["10.32.2.0/27",   "10.32.2.32/27",  "10.32.2.64/27"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "ci.engineering.internal"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost("${var.cidr}", 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "${var.cidr}"
  vpc_id  = "${module.vpc.id}"
  name    = "Engineering"
}

module "jenkins" {
  source                 = "./jenkins"

  internal_subnets       = "${module.vpc.internal_subnets}"
  vpc_id                 = "${module.vpc.id}"
  sg_jenkins             = "${module.security_groups.jenkins_master}"
  sg_jenkins_efs         = "${module.security_groups.jenkins_master_efs}"
  sg_internal_elb        = "${module.security_groups.jenkins_master_elb}"
  region                 = "us-west-2"
}

module "peering_prod_infra" {
  source                    = "../modules/peering"

  vpc_id                    = "${module.vpc.id}"
  external_rtb_id           = "${module.vpc.external_rtb_id}"
  raw_route_tables_id       = "${module.vpc.raw_route_tables_id}"

  tools_account_number      = "${data.terraform_remote_state.prod_infrastructure.account_number}"
  tools_cidr                = "${data.terraform_remote_state.prod_infrastructure.west_cidr}"
  tools_vpc_id              = "${data.terraform_remote_state.prod_infrastructure.west_vpc_id}"
}

resource "aws_vpc_peering_connection" "peer" {
  provider      = "aws.local"

  peer_owner_id = "961082193871"
  peer_vpc_id   = "vpc-10c9f477"
  vpc_id        = "${module.vpc.id}"

  accepter {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "peer_internal_1" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[0]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_2" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[1]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_3" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.raw_route_tables_id[2]}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.local"
  route_table_id            = "${module.vpc.external_rtb_id}"
  destination_cidr_block    = "10.31.0.0/23"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
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
