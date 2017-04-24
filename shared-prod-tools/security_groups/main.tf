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

# Security Group for Tools ECS Instances
resource "aws_security_group" "tools" {
  name        = "engineering-tools"
  description = "Centralized SG for the Shared Production Tools ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  tags {
    Name        = "Shared Production Tools"
  }
}

resource "aws_security_group" "vpn" {
  name        = "${format("%s-vpn", var.name)}"
  description = "Pritunl OpenVPN Server"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 10000
    to_port     = 20000
    protocol    = "udp"
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
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "Pritunl OpenVPN Server"
  }
}

resource "aws_security_group" "internal_elb" {
  name        = "${format("%s-internal_elb", var.name)}"
  description = "Allows ELB Access Access"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.cidr}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.cidr}"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${format("%s internal elb", var.name)}"
  }
}

resource "aws_security_group" "consul-master" {
  name        = "${format("%s-consul-master", var.name)}"
  description = "Allows Consul Access from the ECS Cluster."
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
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
    Name        = "Consul Master"
  }
}

resource "aws_security_group" "internal_ssh" {
  name        = "${format("%s-internal-ssh", var.name)}"
  description = "Allows ssh from inside the VPN/VPC"
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

resource "aws_security_group_rule" "allow_elb" {
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.internal_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_out" {
  type = "egress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.tools.id}"
}

// Internal SSH allows ssh connections from the external ssh security group.
output "internal_ssh" {
  value = "${aws_security_group.internal_ssh.id}"
}

// Internal ELB allows traffic from the internal subnets.
output "internal_elb" {
  value = "${aws_security_group.internal_elb.id}"
}

output "tools-ecs-cluster" {
  value = "${aws_security_group.tools.id}"
}

output "consul-master" {
  value = "${aws_security_group.consul-master.id}"
}

output "vpn" {
  value = "${aws_security_group.vpn.id}"
}