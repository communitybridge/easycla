variable "tools_account_number" {
  description = "The Account Number of the Account we need to send the Peering Request to"
}

variable "tools_vpc_id" {
  description = "The VPC ID on the other side of the Request, the VPC we want to peer to."
}

variable "raw_route_tables_id" {
  type = "list"
  description = "The list of the route tables id of the Internal Subnets."
}

variable "external_rtb_id" {
  description = "The Route Table ID of the External Route Table."
}

variable "tools_cidr" {
  description = "The CIDR of the VPC we are peering to."
}

variable "vpc_id" {
  description = "The VPC ID on our side of the connection. The VPC we are peering FROM."
}

resource "aws_vpc_peering_connection" "peer" {
  provider      = "aws.local"

  peer_owner_id = "${var.tools_account_number}"
  peer_vpc_id   = "${var.tools_vpc_id}"
  vpc_id        = "${var.vpc_id}"

  accepter {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "peer_internal_1" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[0]}"
  destination_cidr_block    = "${var.tools_cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_2" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[1]}"
  destination_cidr_block    = "${var.tools_cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_internal_3" {
  provider                  = "aws.local"
  route_table_id            = "${var.raw_route_tables_id[2]}"
  destination_cidr_block    = "${var.tools_cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}

resource "aws_route" "peer_external" {
  provider                  = "aws.local"
  route_table_id            = "${var.external_rtb_id}"
  destination_cidr_block    = "${var.tools_cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
}
