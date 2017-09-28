variable "key_name" {}

variable "external_subnets" {
  type = "list"
}

variable "vpc_id" {}

variable "cidr" {}

variable "route_tables" {
  type = "list"
}

resource "aws_security_group" "it-vpn-tunnel" {
  provider    = "aws.local"
  name        = "cert-managed-vpn-tunnel"
  description = "Cert Managed VPN Tunnel"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.cidr}"]
  }

  # Certification firewalls.
  #  - 54.213.24.55     (us-west-2-lfc-fw-1.dmz)
  #  - 174.143.200.141  (dfw-lfc-fw.dmz)
  #  - 119.9.92.26      (hkg-lfc-fw.dmz)
  #  - 54.193.90.43     (us-west-1-lfc-fw.dmz)
  #  - 54.172.59.54     (us-east-1-lfc-fw.dmz)
  #  - 54.179.151.230   (ap-southeast-1-lfc-fw.dmz)
  #  - 52.193.246.225   (ap-northeast-1-lfc-fw.dmz)
  ingress {
    from_port   = 1194
    to_port     = 1194
    protocol    = "udp"
    cidr_blocks = [
        "54.213.24.55/32",
        "174.143.200.141/32",
        "119.9.92.26/32",
        "54.193.90.43/32",
        "54.172.59.54/32",
        "54.179.151.230/32",
        "52.193.246.225/32"
    ]
  }

  # Rene's managment server in Cert network.
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["52.64.150.59/32", "13.210.89.242"]
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
    Name        = "Cert Managed VPN Tunnel"
  }
}

resource "aws_instance" "it-managed-vpn-tunnel" {
  ami                     = "ami-d2c924b2" # Centos 7
  instance_type           = "t2.medium"
  key_name                = "${var.key_name}"
  count                   = "1"
  source_dest_check       = "false"

  root_block_device {
    volume_size           = 50
  }

  disable_api_termination = "true"
  vpc_security_group_ids  = ["${aws_security_group.it-vpn-tunnel.id}"]
  subnet_id               = "${var.external_subnets[0]}"

  tags {
    Name                  = "Cert Managed VPN Tunnel"
  }
}
