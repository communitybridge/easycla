variable "name" {
  description = "The name will be used to prefix and tag the resources, e.g myredis"
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

variable "instance_type" {
  description = "The type of instances that the Redis cluster will be running on"
  default     = "cache.t2.small"
}

variable "instance_count" {
  description = "How many nodes will be provisioned in the Redis cluster"
  default     = 1
}

variable "preferred_backup_window" {
  description = "The time window on which backups will be made (HH:mm-HH:mm)"
  default     = "07:00-09:00"
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
  description = "Route53 record name for the Redis Cluster, defaults to the database name if not set"
  default     = ""
}

variable "port" {
  description = "The port at which redis listens for incoming connections"
  default     = 6379
}

variable "parameter_group_name" {
  description = "Redis Parameter Group Name"
  default     = "default.redis3.2"
}

variable "version" {
  description = "Redis Version"
  default     = "3.2.4"
}

variable "az_mode" {
  description = "Redis AZ Mode (single-az or cross-az)"
  default     = "single-az"
}

resource "aws_security_group" "main" {
  name        = "${var.name}-redis-node"
  description = "Allows traffic to redis from other security groups"
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
    Name        = "Redis node (${var.name})"
    Environment = "${var.environment}"
  }
}

resource "aws_elasticache_subnet_group" "main" {
  name        = "${var.name}"
  description = "Redis subnet group"
  subnet_ids  = ["${var.subnet_ids}"]
}

resource "aws_elasticache_cluster" "main" {
  cluster_id              = "${var.name}"
  engine                  = "redis"
  engine_version          = "${var.version}"
  node_type               = "${var.instance_type}"
  port                    = "${var.port}"
  num_cache_nodes         = "${var.instance_count}"
  parameter_group_name    = "${var.parameter_group_name}"
  maintenance_window      = "${var.preferred_maintenance_window}"
  az_mode                 = "${var.az_mode}"
  security_group_ids      = ["${aws_security_group.main.id}"]
  subnet_group_name       = "${aws_elasticache_subnet_group.main.name}"

  tags {
    Name        = "${var.name}"
    Environment = "${var.environment}"
    Team        = "${var.team}"
  }
}

// The cluster identifier.
output "id" {
  value = "${aws_elasticache_cluster.main.id}"
}

output "endpoint" {
  value = "${aws_elasticache_cluster.main.endpoint}"
}

output "port" {
  value = "${aws_elasticache_cluster.main.port}"
}
