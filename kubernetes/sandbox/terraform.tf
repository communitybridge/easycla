variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

//terraform {
//  backend "consul" {
//    address = "consul.service.production.consul:8500"
//    path    = "terraform/infrastructure"
//  }
//}

provider "aws" {
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

module "iam_user" "kops" {
  source = "../modules/kops_iam_user"
}