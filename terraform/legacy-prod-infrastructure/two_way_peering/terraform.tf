variable "vpc_id" {}

variable "internal_subnets" {
  type = "list"
}

variable "cidr" {}

variable "external_rtb_id" {}

variable "raw_route_tables_id" {
  type = "list"
}

variable "s3_state_path" {}

data "terraform_remote_state" "engineering" {
  backend = "s3"
  config {
    bucket = "lfe-terraform-states"
    key = "${var.s3_state_path}"
    region = "us-west-2"
    access_key = "AKIAJZSPEM5HOEMPP67Q"
    secret_key = "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"
  }
}

resource "aws_vpc_peering_connection" "peer" {
  provider      = "aws.western"
  peer_owner_id = "433610389961"
  peer_vpc_id   = "${data.terraform_remote_state.engineering.vpc_id}"
  vpc_id        = "${var.vpc_id}"
  auto_accept   = true

  accepter {
    allow_remote_vpc_dns_resolution = true
  }

  requester {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "peer_internal" {
  provider                  = "aws.western"
  count                     = "${length(compact(var.internal_subnets))}"
  route_table_id            = "${var.raw_route_tables_id[count.index]}"
  destination_cidr_block    = "${data.terraform_remote_state.engineering.cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.western"
  route_table_id            = "${var.external_rtb_id}"
  destination_cidr_block    = "${data.terraform_remote_state.engineering.cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "internal" {
  provider               = "aws.western"
  count                  = "${length(compact(data.terraform_remote_state.engineering.internal_subnets))}"
  route_table_id         = "${data.terraform_remote_state.engineering.raw_route_tables_id[count.index]}"
  destination_cidr_block = "${var.cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

output "vpx" {
  value = "${aws_vpc_peering_connection.peer.id}"
}