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
  provider    = "aws.local"
  name        = "engineering-tools"
  description = "Centralized SG for the Shared Production Tools ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    security_groups = ["${aws_security_group.internal_elb.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    security_groups = ["${aws_security_group.vpn.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    security_groups = ["${aws_security_group.bind.id}"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    security_groups = ["433610389961/sg-0d2f4376"] # ENG JENKINS MASTER
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    security_groups = ["433610389961/sg-493b5732"] # ENG JENKINS SLAVES
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Shared Production Tools"
  }
}

resource "aws_security_group" "vpn" {
  provider    = "aws.local"
  name        = "${format("%s-vpn", var.name)}"
  description = "Pritunl OpenVPN Server"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["50.188.159.47/32"]
  }

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

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
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
    Name        = "Pritunl OpenVPN Server"
  }
}

resource "aws_security_group" "internal_elb" {
  provider    = "aws.local"
  name        = "${format("%s-internal_elb", var.name)}"
  description = "Allows ELB Access Access"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
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
    Name        = "${format("%s internal elb", var.name)}"
  }
}

resource "aws_security_group" "bind" {
  provider    = "aws.local"
  name        = "${format("%s-bind-servers", var.name)}"
  description = "Allows Access to DNS Servers"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 53
    to_port     = 53
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 53
    to_port     = 53
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    security_groups = ["${aws_security_group.vpn.id}"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
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
    Name        = "${format("%s bind servers", var.name)}"
  }
}

resource "aws_security_group" "production-tools-efs" {
  provider    = "aws.local"
  name        = "production-tools-efs"
  description = "Allows Internal EFS use from Production Tools"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = ["${aws_security_group.tools.id}"]
  }

  egress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "Production Tools EFS"
  }
}

resource "aws_security_group" "vpn-link" {
  provider    = "aws.local"
  name        = "${format("%s-vpn-link", var.name)}"
  description = "Pritunl OpenVPN Link"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 4500
    to_port     = 4500
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 500
    to_port     = 500
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl Link"
  }
}

resource "aws_security_group" "ghe" {
  provider    = "aws.local"
  name        = "github-enterprise"
  description = "GitHub Enterprise Appliance"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    security_groups = ["${aws_security_group.ghe_elb.id}"]
  }

  ingress {
    from_port = 25
    to_port = 25
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    security_groups = ["${aws_security_group.ghe_elb.id}"]
  }

  ingress {
    from_port = 122
    to_port = 122
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 123
    to_port = 123
    protocol = "udp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 161
    to_port = 161
    protocol = "udp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 8443
    to_port = 8443
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    security_groups = ["${aws_security_group.ghe_elb.id}"]
  }

  ingress {
    from_port = 1194
    to_port = 1194
    protocol = "udp"
    self = true
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
    Name        = "GitHub Enterprise Appliance"
  }
}

resource "aws_security_group" "ghe_elb" {
  provider    = "aws.local"
  name        = "github-enterprise-elb"
  description = "GitHub Enterprise Appliance"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
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
    Name        = "GitHub Enterprise Appliance"
  }
}

// Internal ELB allows traffic from the internal subnets.
output "internal_elb" {
  value = "${aws_security_group.internal_elb.id}"
}

output "tools-ecs-cluster" {
  value = "${aws_security_group.tools.id}"
}

output "ghe-elb" {
  value = "${aws_security_group.ghe_elb.id}"
}

output "vpn" {
  value = "${aws_security_group.vpn.id}"
}

output "vpn_link" {
  value = "${aws_security_group.vpn-link.id}"
}

output "bind" {
  value = "${aws_security_group.bind.id}"
}

output "efs" {
  value = "${aws_security_group.production-tools-efs.id}"
}

output "ghe" {
  value = "${aws_security_group.ghe.id}"
}