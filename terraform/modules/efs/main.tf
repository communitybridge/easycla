variable "name" {
  description = "The cluster name, e.g cdn"
}

variable "display_name" {
  description = "The cluster name, e.g cdn"
}

variable "vpc_id" {
  description = "VPC ID"
}

variable "vpc_cidr" {
  description = "VPC CIDR"
}

variable "subnet_ids" {
  description = "List of subnet IDs"
  type        = "list"
}

variable "security_group" {
  description = "Security group to use for the ECS Cluster"
}

# Security Group for EFS
resource "aws_security_group" "efs" {
  provider           = "aws.local"
  name        = "${var.name}-efs"
  description = "Allows Internal EFS use from inside the VPC"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = ["${var.security_group}"]
  }

  egress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = ["${var.vpc_cidr}"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "${var.display_name}"
  }
}

# Creating EFS for Tools Storage
resource "aws_efs_file_system" "efs" {
  provider           = "aws.local"
  creation_token = "${var.name}-efs"
  tags {
    Name = "${var.display_name}"
  }
}

resource "aws_efs_mount_target" "efs_mount_1" {
  provider           = "aws.local"
  file_system_id = "${aws_efs_file_system.efs.id}"
  subnet_id = "${var.subnet_ids[0]}"
  security_groups = ["${aws_security_group.efs.id}"]
}

resource "aws_efs_mount_target" "efs_mount_2" {
  provider           = "aws.local"
  file_system_id = "${aws_efs_file_system.efs.id}"
  subnet_id = "${var.subnet_ids[1]}"
  security_groups = ["${aws_security_group.efs.id}"]
}

resource "aws_efs_mount_target" "efs_mount_3" {
  provider           = "aws.local"
  file_system_id = "${aws_efs_file_system.efs.id}"
  subnet_id = "${var.subnet_ids[2]}"
  security_groups = ["${aws_security_group.efs.id}"]
}

// The EFS ID
output "id" {
  value = "${aws_efs_file_system.efs.id}"
}

// The SG ID
output "sg" {
  value = "${aws_security_group.efs.id}"
}

