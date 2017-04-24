data "terraform_remote_state" "shared-prod-tools" {
  backend = "s3"
  config {
    bucket = "lfe-terraform-states"
    key = "shared-prod-tools/terraform.tfstate"
    region = "us-west-2"
    access_key = "AKIAJZSPEM5HOEMPP67Q"
    secret_key = "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"
  }
}

resource "aws_route" "internal" {
  count                  = "${length(compact(module.vpc.internal_subnets))}"
  route_table_id         = "${module.vpc.raw_route_tables_id[count.index]}"
  destination_cidr_block = "${data.terraform_remote_state.shared-prod-tools.cidr}"
  vpc_peering_connection_id = "${data.terraform_remote_state.shared-prod-tools.vpx_to_engineering}"
}