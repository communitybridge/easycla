variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

data "terraform_remote_state" "engineering" {
  backend = "s3"
  config {
    bucket = "lfe-terraform-states"
    key = "general/terraform.tfstate"
    region = "us-west-2"
  }
}

resource "aws_s3_bucket" "engineering-database-backups" {
  bucket = "engineering-database-backups"
  acl = "private"

  tags {
    Name = "Engineering Database Backups"
    Environment = "Production"
  }
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/files/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "engineering-preprods"
    newrelic_hostname= "Engineering Preprods - ECS Container Instance"
    newrelic_labels  = "Environment:Preprod;Team:Engineering;Datacenter:AWS"
    newrelic_key     = "951db34ebed364ea663002571b63db5d3f827758"
  }
}

provider "aws" {
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

# Security Group for EFS
resource "aws_security_group" "preprods" {
  name        = "engineering-preprods"
  description = "Centralized SG for all the Preprods using the Tools ECS Cluster"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  tags {
    Name        = "Preprods: ECS Cluster"
  }
}

resource "aws_security_group_rule" "allow_elb" {
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${data.terraform_remote_state.engineering.sg_external_elb}"

  security_group_id = "${aws_security_group.preprods.id}"
}

resource "aws_security_group_rule" "allow_out" {
  type = "egress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.preprods.id}"
}

module "engineering-preprods-ecs-cluster" {
  source               = "../../modules/ecs-cluster"
  environment          = "Preprod"
  team                 = "Engineering"
  name                 = "engineering-preprods"
  vpc_id               = "${data.terraform_remote_state.engineering.vpc_id}"
  subnet_ids           = "${data.terraform_remote_state.engineering.external_subnets}"
  key_name             = "engineering-preprods"
  iam_instance_profile = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region               = "${data.terraform_remote_state.engineering.region}"
  availability_zones   = "${data.terraform_remote_state.engineering.availability_zones}"
  instance_type        = "t2.medium"
  security_group       = "${aws_security_group.preprods.id}"
  instance_ebs_optimized = false
  desired_capacity     = "3"
  min_size             = "3"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
}

module "rds-cluster" {
  source               = "../../modules/rds-cluster"
  master_username      = "lfengineering"
  name                 = "engineering-preprods"
  master_password      = "buanCAWwwAGxUyoU2Fai"
  zone_id              = "${data.terraform_remote_state.engineering.zone_id}"
  availability_zones   = "${data.terraform_remote_state.engineering.availability_zones}"
  vpc_id               = "${data.terraform_remote_state.engineering.vpc_id}"
  subnet_ids           = "${data.terraform_remote_state.engineering.internal_subnets}"
  environment          = "Preprod"
  team                 = "Engineering"
  security_groups      = ["${module.engineering-preprods-ecs-cluster.security_group_id}"]
  engine               = "mariadb"
  engine_version       = "10.1.19"
  parameter_group_name = "engineering"
  dns_name             = "rds"
}

module "redis-cluster" {
  source          = "../../modules/redis-cluster"
  name            = "engineering-preprods"
  version         = "3.2.4"
  instance_type   = "cache.t2.medium"
  instance_count  = "1"
  environment     = "Preprod"
  team            = "Engineering"
  zone_id         = "${data.terraform_remote_state.engineering.zone_id}"
  security_groups = ["${module.engineering-preprods-ecs-cluster.security_group_id}"]
  subnet_ids      = "${data.terraform_remote_state.engineering.internal_subnets}"
  vpc_id          = "${data.terraform_remote_state.engineering.vpc_id}"
  dns_name        = "redis"
}