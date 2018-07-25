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

resource "aws_security_group" "engineering-production-elb" {
  provider    = "aws.local"
  name        = "production-elb"
  description = "Public ELB used for Production"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 80
    to_port         = 80
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
    from_port = 9990
    protocol = "tcp"
    to_port = 9990
    cidr_blocks     = ["0.0.0.0/0"]
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
    Name        = "Sandboxes ELB"
  }
}

resource "aws_security_group" "engineering-production" {
  provider    = "aws.local"
  name        = "production-ecs"
  description = "Production ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.engineering-production-elb.id}"]
  }

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    self            = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Sandboxes ECS Cluster"
  }
}

output "engineering_production" {
  value = "${aws_security_group.engineering-production.id}"
}

output "engineering_production_elb" {
  value = "${aws_security_group.engineering-production-elb.id}"
}