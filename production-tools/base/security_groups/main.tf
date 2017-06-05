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
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    security_groups = ["${aws_security_group.tools.id}"]
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

resource "aws_security_group" "internal_ssh" {
  provider    = "aws.local"
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
  provider    = "aws.local"
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.internal_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_vpn" {
  provider    = "aws.local"
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.vpn.id}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_itself" {
  provider    = "aws.local"
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.tools.id}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_bind" {
  provider    = "aws.local"
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.bind.id}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_out" {
  provider    = "aws.local"
  type = "egress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.tools.id}"
}

//resource "aws_security_group_rule" "allow_cinco_prod" {
//  provider    = "aws.local"
//  type = "ingress"
//  from_port = 0
//  to_port = 0
//  protocol = "-1"
//  source_security_group_id = "643009352547/sg-ffcc8d84"
//
//  security_group_id = "${aws_security_group.tools.id}"
//}
//
//resource "aws_security_group_rule" "allow_pmc_prod" {
//  provider    = "aws.local"
//  type = "ingress"
//  from_port = -1
//  to_port = -1
//  protocol = "-1"
//  source_security_group_id = "643009352547/sg-a64c08dd"
//
//  security_group_id = "${aws_security_group.tools.id}"
//}

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