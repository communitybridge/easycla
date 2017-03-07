variable "name" {
  description = "The name will be used to prefix and tag the resources, e.g mydb"
}

variable "environment" {
  description = "The environment tag, e.g prod"
}

variable "team" {
  description = "Team tag, e.g Engineering"
}

variable "vpc_id" {
  description = "The VPC ID to use"
}

variable "zone_id" {
  description = "The Route53 Zone ID where the DNS record will be created"
}

variable "security_groups" {
  description = "A list of security group IDs"
  type = "list"
}

variable "subnet_ids" {
  description = "A list of subnet IDs"
  type = "list"
}

variable "availability_zones" {
  description = "A list of availability zones"
  type = "list"
}

variable "master_username" {
  description = "The master user username"
}

variable "master_password" {
  description = "The master user password"
}

variable "instance_type" {
  description = "The type of instances that the RDS cluster will be running on"
  default     = "db.r3.large"
}

variable "instance_count" {
  description = "How many instances will be provisioned in the RDS cluster"
  default     = 1
}

variable "preferred_backup_window" {
  description = "The time window on which backups will be made (HH:mm-HH:mm)"
  default     = "07:00-09:00"
}

variable "backup_retention_period" {
  description = "The backup retention period"
  default     = 5
}

variable "preferred_maintenance_window" {
  description = "The time window on which maintenance will be made (ddd:hh24:mi-ddd:hh24:mi)"
  default     = "Mon:00:00-Mon:03:00"
}

variable "publicly_accessible" {
  description = "When set to true the RDS cluster can be reached from outside the VPC"
  default     = false
}

variable "dns_name" {
  description = "Route53 record name for the RDS database, defaults to the database name if not set"
  default     = ""
}

variable "port" {
  description = "The port at which the database listens for incoming connections"
  default     = 3306
}

variable "multi_az" {
  description = "Should RDS be used in a Multi-AZ fashion"
  default     = false
}

variable "engine" {
  description = "The Database Engine to be used"
  default     = "mariadb"
}

variable "engine_version" {
  description = "Database Engine Version"
  default     = "10.1.19"
}

variable "parameter_group_name" {
  description = "RDS Database Parameter Group Name"
  default     = "default.mariadb10.1"
}

variable "allocated_storage" {
  description = "Allocated Storage to the RDS Instance"
  default     = 100
}

resource "aws_security_group" "main" {
  name        = "${var.name}-rds-cluster"
  description = "Allows traffic to rds from other security groups"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = "${var.port}"
    to_port         = "${var.port}"
    protocol        = "TCP"
    security_groups = ["${var.security_groups}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "RDS cluster (${var.name})"
    Environment = "${var.environment}"
  }
}

resource "aws_db_subnet_group" "main" {
  name        = "${var.name}"
  description = "RDS cluster subnet group"
  subnet_ids  = ["${var.subnet_ids}"]
}

resource "aws_db_instance" "main" {
  allocated_storage       = "${var.allocated_storage}"
  engine                  = "${var.engine}"
  engine_version          = "${var.engine_version}"
  instance_class          = "${var.instance_type}"
  identifier              = "${var.name}"
  username                = "${var.master_username}"
  password                = "${var.master_password}"
  db_subnet_group_name    = "${aws_db_subnet_group.main.id}"
  parameter_group_name    = "${var.parameter_group_name}"
  backup_retention_period = "${var.backup_retention_period}"
  maintenance_window      = "${var.preferred_maintenance_window}"
  vpc_security_group_ids  = ["${aws_security_group.main.id}"]
  backup_window           = "${var.preferred_backup_window}"
  port                    = "${var.port}"
  multi_az                = "${var.multi_az}"
  publicly_accessible     = "${var.publicly_accessible}"

  tags {
    Name        = "${var.name}"
    Environment = "${var.environment}"
    Team        = "${var.team}"
  }
}

// The cluster identifier.
output "id" {
  value = "${aws_db_instance.main.id}"
}

output "endpoint" {
  value = "${aws_db_instance.main.endpoint}"
}

output "port" {
  value = "${aws_db_instance.main.port}"
}
