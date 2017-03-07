variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

provider "aws" {
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

data "terraform_remote_state" "engineering" {
  backend = "s3"
  config {
    bucket = "lfe-terraform-states"
    key = "general/terraform.tfstate"
    region = "us-west-2"
  }
}

resource "aws_instance" "openvpn" {
  ami = "ami-a6d773c6"
  instance_type = "m3.medium"
  key_name = "engineering-vpn"
  subnet_id = "${data.terraform_remote_state.engineering.external_subnets[0]}"
  vpc_security_group_ids = ["${data.terraform_remote_state.engineering.sg_vpn}"]

  tags {
    Name = "OpenVPN Server"
  }

}

resource "aws_eip" "openvpn" {
  instance = "${aws_instance.openvpn.id}"
  vpc      = true
}

resource "aws_route53_record" "vpn" {
  zone_id = "Z2MDT77FL23F9B"
  name = "vpn.engineering.tux.rocks."
  type = "A"
  ttl = "300"
  records = ["${aws_eip.openvpn.public_ip}"]
}