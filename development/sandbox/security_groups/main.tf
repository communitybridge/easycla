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

resource "aws_security_group" "engineering-sandboxes-elb" {
  provider    = "aws.local"
  name        = "sandboxes-elb"
  description = "Public ELB used for Sandboxes"
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

resource "aws_security_group" "engineering-sandboxes" {
  provider    = "aws.local"
  name        = "sandboxes-ecs"
  description = "Sandboxes ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.engineering-sandboxes-elb.id}"]
  }

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.engineering-sandboxes-redis.id}"]
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

resource "aws_security_group" "engineering-sandboxes-redis" {
  provider    = "aws.local"
  name        = "sandboxes-redis"
  description = "Sandboxes Redis"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    cidr_blocks     = ["${var.cidr}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Sandboxes Redis Server"
  }
}

output "engineering_sandboxes" {
  value = "${aws_security_group.engineering-sandboxes.id}"
}

output "engineering_sandboxes_elb" {
  value = "${aws_security_group.engineering-sandboxes-elb.id}"
}

output "engineering_sandboxes_redis" {
  value = "${aws_security_group.engineering-sandboxes-redis.id}"
}