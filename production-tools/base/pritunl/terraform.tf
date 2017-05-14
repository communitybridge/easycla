variable "external_subnets" {
  description = "External VPC Subnets"
  type = "list"
}

variable "vpn_sg" {
  description = "Security Group for the VPN Server"
}

variable "region_identifier" {}

data "aws_ami" "pritunl" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["pritunl-2*"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "pritunl" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.pritunl.id}"
  source_dest_check      = false
  instance_type          = "t2.small"
  subnet_id              = "${element(var.external_subnets, 0)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.vpn_sg}"]
  monitoring             = true

  tags {
    Name        = "Pritunl - VPN Instance"
    Team        = "Engineering"
    Environment = "Production"
  }
}

resource "aws_eip" "pritunl" {
  provider = "aws.local"
  instance = "${aws_instance.pritunl.id}"
  vpc      = true
}

resource "aws_route53_record" "pritunl" {
  zone_id = "Z2MDT77FL23F9B"
  name    = "${var.region_identifier}.vpn"
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.pritunl.public_ip}"]
}
