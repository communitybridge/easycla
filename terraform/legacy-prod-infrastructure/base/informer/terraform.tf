variable "subnets" {
  type = "list"
}

resource "aws_rds_cluster" "informer" {
  cluster_identifier      = "informer-5"
  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "informer"
  master_username         = "lfengineering"
  master_password         = "Es}Fpo2zH&7eEF3kw)GePoXz"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
}

resource "aws_db_subnet_group" "informer" {
  name       = "main"
  subnet_ids = ["${var.subnets}"]

  tags {
    Name = "Informer Subnet Group"
  }
}

resource "aws_elasticache_cluster" "informer" {
  cluster_id           = "informer-5"
  engine               = "redis"
  node_type            = "cache.t2.small"
  port                 = 11211
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
}

resource "aws_elasticache_subnet_group" "informer" {
  name       = "informer-subnet-group"
  subnet_ids = ["${var.subnets}"]
}