data "terraform_remote_state" "engineering" {
  backend = "s3"
  config {
    bucket = "lfe-terraform-states"
    key = "general/terraform.tfstate"
    region = "us-west-2"
    access_key = "AKIAJZSPEM5HOEMPP67Q"
    secret_key = "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"
  }
}

resource "aws_vpc_peering_connection" "engineering_sandbox" {
  peer_owner_id = "433610389961"
  peer_vpc_id   = "${data.terraform_remote_state.engineering.vpc_id}"
  vpc_id        = "${module.vpc.id}"
  auto_accept   = true

  accepter {
    allow_remote_vpc_dns_resolution = true
  }

  requester {
    allow_remote_vpc_dns_resolution = true
  }
}

resource "aws_route" "engineering_sandbox_vpx" {
  count                  = "${length(compact(module.vpc.internal_subnets))}"
  route_table_id         = "${module.vpc.raw_route_tables_id[count.index]}"
  destination_cidr_block = "${data.terraform_remote_state.engineering.cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.engineering_sandbox.id}"
}

resource "aws_route" "engineering_sandbox_vpx_external" {
  route_table_id         = "${module.vpc.external_rtb_id}"
  destination_cidr_block = "${data.terraform_remote_state.engineering.cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection.engineering_sandbox.id}"
}

resource "aws_route" "coreit_vpx" {
  count                  = "${length(compact(module.vpc.internal_subnets))}"
  route_table_id         = "${module.vpc.raw_route_tables_id[count.index]}"
  destination_cidr_block = "172.17.0.0/18"
  vpc_peering_connection_id = "pcx-79ac2510"
}

resource "aws_route" "coreit_vpx_external" {
  route_table_id         = "${module.vpc.external_rtb_id}"
  destination_cidr_block = "172.17.0.0/18"
  vpc_peering_connection_id = "pcx-79ac2510"
}

output "vpx_to_engineering" {
  value = "${aws_vpc_peering_connection.engineering_sandbox.id}"
}