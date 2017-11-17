/**
 * Creates basic security groups to be used by instances and ELBs.
 */

variable "vpc_id" {
  description = "The VPC ID"
}

resource "aws_security_group" "pritunl-elb" {
  provider    = "aws.local"
  name        = "pritunl-elb"
  description = "ELB in front of Pritunl Nodes"
  vpc_id      = "${var.vpc_id}"

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
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl Load Balancer"
  }
}

resource "aws_security_group" "pritunl-node" {
  provider    = "aws.local"
  name        = "pritunl-node"
  description = "Pritunl Server Nodes"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    security_groups = ["${aws_security_group.pritunl-elb.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl Server Nodes"
  }
}

resource "aws_security_group" "pritunl-mongodb" {
  provider    = "aws.local"
  name        = "pritunl-mongodb"
  description = "Pritunl MongoDB Storage"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 27017
    to_port   = 27017
    protocol  = "tcp"
    security_groups = ["${aws_security_group.pritunl-node.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl MongoDB Storage"
  }
}

output "pritunl_elb" {
  value = "${aws_security_group.pritunl-elb.id}"
}

output "pritunl_node" {
  value = "${aws_security_group.pritunl-node.id}"
}

output "pritunl_mongodb" {
  value = "${aws_security_group.pritunl-mongodb.id}"
}