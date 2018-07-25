variable "project_cidr" {}

variable "raw_route_tables_id" {
  type = "list"
}

variable "external_rtb_id" {}

variable "peering_id" {}

resource "aws_route" "peer_internal_1" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[0]}"
  destination_cidr_block    = "${var.project_cidr}"
  vpc_peering_connection_id = "${var.peering_id}"
}

resource "aws_route" "peer_internal_2" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[1]}"
  destination_cidr_block    = "${var.project_cidr}"
  vpc_peering_connection_id = "${var.peering_id}"
}

resource "aws_route" "peer_internal_3" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[2]}"
  destination_cidr_block    = "${var.project_cidr}"
  vpc_peering_connection_id = "${var.peering_id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.local"
  route_table_id            = "${var.external_rtb_id}"
  destination_cidr_block    = "${var.project_cidr}"
  vpc_peering_connection_id = "${var.peering_id}"
}