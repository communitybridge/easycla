variable "key_name" {}

variable "internal_subnets" {
  type = "list"
}

variable "vpc_id" {}

variable "cidr" {}

variable "route_tables" {
  type = "list"
}

resource "aws_security_group" "it-vpn-tunnel" {
  provider    = "aws.local"
  name        = "it-managed-vpn-tunnel"
  description = "IT Managed VPN Tunnel"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.cidr}"]
  }

  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [
      "10.32.0.0/12",
      "98.234.51.222/32",
      "103.12.132.88/32"
    ]
  }

  ingress {
    from_port   = 122
    to_port     = 122
    protocol    = "tcp"
    cidr_blocks = [
      "10.32.0.0/12",
      "98.234.51.222/32",
      "103.12.132.88/32"
    ]
  }

  ingress {
    from_port   = 8443
    to_port     = 8443
    protocol    = "tcp"
    cidr_blocks = [
      "10.32.0.0/12",
      "98.234.51.222/32",
      "103.12.132.88/32"
    ]
  }

  ingress {
    from_port   = 636
    to_port     = 636
    protocol    = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "IT Managed VPN Tunnel"
  }
}

resource "aws_instance" "it-managed-vpn-tunnel" {
  ami                     = "ami-8652b5fe" # AMI for it-managed vpn
  instance_type           = "t2.medium"
  key_name                = "${var.key_name}"
  count                   = "1"
  source_dest_check       = "false"

  root_block_device {
    volume_size           = 50
  }

  disable_api_termination = "true"
  vpc_security_group_ids  = ["${aws_security_group.it-vpn-tunnel.id}"]
  subnet_id               = "${var.internal_subnets[0]}"

  tags {
    Name                  = "IT Managed VPN Tunnel"
  }
}

resource "aws_eip" "it-managed-vpn-tunnel" {
  vpc                       = true
  instance                  = "${aws_instance.it-managed-vpn-tunnel.id}"
}

resource "aws_route" "route_it_1" {
  count                     = "${length(var.internal_subnets)}"
  provider                  = "aws.local"
  route_table_id            = "${var.route_tables[count.index]}"
  destination_cidr_block    = "10.30.0.0/16"
  instance_id               = "${aws_instance.it-managed-vpn-tunnel.id}"
}

resource "aws_route" "route_it_2" {
  count                     = "${length(var.internal_subnets)}"
  provider                  = "aws.local"
  route_table_id            = "${var.route_tables[count.index]}"
  destination_cidr_block    = "172.30.100.128/25"
  instance_id               = "${aws_instance.it-managed-vpn-tunnel.id}"
}