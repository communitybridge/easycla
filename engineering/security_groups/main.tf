/**
 * Creates basic security groups to be used by instances and ELBs.
 */

variable "name" {
  description = "The name of the security groups serves as a prefix, e.g stack"
}

variable "vpc_id" {
  description = "The VPC ID"
}

variable "cidr" {
  description = "The cidr block to use for internal security groups"
}

resource "aws_security_group" "vpn" {
  name        = "engineering-vpn"
  description = "Allows External Access to the OpenVPN EndPoint"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 943
    to_port         = 943
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 1194
    to_port         = 1194
    protocol        = "udp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "engineering-openvpn"
  }
}

resource "aws_security_group" "internal_web_elb" {
  name        = "${format("%s-internal-web-elb", var.name)}"
  vpc_id      = "${var.vpc_id}"
  description = "Allows internal ELB 80/443 traffic"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["${var.cidr}"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["${var.cidr}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${format("%s internal elb", var.name)}"
  }
}

resource "aws_security_group" "external_web_elb" {
  name        = "${format("%s-external-web-elb", var.name)}"
  vpc_id      = "${var.vpc_id}"
  description = "Allows external ELB 80/443 traffic"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${format("%s external elb", var.name)}"
  }
}

resource "aws_security_group" "external_ssh" {
  name        = "${format("%s-external-ssh", var.name)}"
  description = "Allows ssh from the world"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["140.211.169.2/32"]
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
    Name        = "${format("%s external ssh", var.name)}"
  }
}

resource "aws_security_group" "internal_ssh" {
  name        = "${format("%s-internal-ssh", var.name)}"
  description = "Allows ssh from bastion"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "tcp"
    cidr_blocks = ["${var.cidr}"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${format("%s internal ssh", var.name)}"
  }
}

// External SSH allows ssh connections on port 22 from the world.
output "external_ssh" {
  value = "${aws_security_group.external_ssh.id}"
}

// Internal SSH allows ssh connections from the external ssh security group.
output "internal_ssh" {
  value = "${aws_security_group.internal_ssh.id}"
}

// External ELB allows traffic from the world.
output "external_elb" {
  value = "${aws_security_group.external_web_elb.id}"
}

// Internal ELB allows traffic from the internal subnets.
output "internal_elb" {
  value = "${aws_security_group.internal_web_elb.id}"
}

// Internal subnets allows traffic from the VPN.
output "vpn" {
  value = "${aws_security_group.vpn.id}"
}
