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
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
  }
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/files/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "engineering-tools"
    efs_id           = "${aws_efs_file_system.engineering-tools-storage.id}"
    efs_name         = "engineering_tools_storage"
    newrelic_hostname= "Engineering Tools - ECS Container Instance"
    newrelic_labels  = "Environment:Production;Team:Engineering;Datacenter:AWS"
    newrelic_key     = "951db34ebed364ea663002571b63db5d3f827758"
  }
}

provider "aws" {
  region = "us-west-2"
  access_key = "AKIAJ5HCOKTEUSBBBVCQ"
  secret_key = "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"
}

resource "aws_cloudwatch_log_group" "tools" {
  name = "engineering-tools"
}

# Creating EFS for Tools Storage
resource "aws_efs_file_system" "engineering-tools-storage" {
  creation_token = "enginnering-tools-storage"
  tags {
    Name = "Enginnering Tools - Storage"
  }
}

resource "aws_efs_mount_target" "efs_mount_1" {
  file_system_id = "${aws_efs_file_system.engineering-tools-storage.id}"
  subnet_id = "${data.terraform_remote_state.engineering.internal_subnets[0]}"
  security_groups = ["${aws_security_group.sg_internal_efs.id}"]
}

resource "aws_efs_mount_target" "efs_mount_2" {
  file_system_id = "${aws_efs_file_system.engineering-tools-storage.id}"
  subnet_id = "${data.terraform_remote_state.engineering.internal_subnets[1]}"
  security_groups = ["${aws_security_group.sg_internal_efs.id}"]
}

resource "aws_efs_mount_target" "efs_mount_3" {
  file_system_id = "${aws_efs_file_system.engineering-tools-storage.id}"
  subnet_id = "${data.terraform_remote_state.engineering.internal_subnets[2]}"
  security_groups = ["${aws_security_group.sg_internal_efs.id}"]
}

# Security Group for EFS
resource "aws_security_group" "tools" {
  name        = "engineering-tools"
  description = "Centralized SG for all the tools using the Tools ECS Cluster"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  tags {
    Name        = "Tools: ECS Cluster"
  }
}

resource "aws_security_group_rule" "allow_ssh" {
  type = "ingress"
  from_port = 22
  to_port = 22
  protocol = "tcp"
  source_security_group_id = "${data.terraform_remote_state.engineering.sg_internal_ssh}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_vpn" {
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${data.terraform_remote_state.engineering.sg_vpn}"

  security_group_id = "${aws_security_group.tools.id}"
}

resource "aws_security_group_rule" "allow_out" {
  type = "egress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.tools.id}"
}

# Creating ECS-Cluster
module "engineering-tools-ecs-cluster" {
  source               = "../../modules/ecs-cluster"
  environment          = "Production"
  team                 = "Engineering"
  name                 = "engineering-tools"
  vpc_id               = "${data.terraform_remote_state.engineering.vpc_id}"
  subnet_ids           = "${data.terraform_remote_state.engineering.internal_subnets}"
  key_name             = "engineering-tools"
  iam_instance_profile = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region               = "us-west-2"
  availability_zones   = "${data.terraform_remote_state.engineering.availability_zones}"
  instance_type        = "t2.medium"
  security_group      = "${aws_security_group.tools.id}"
  instance_ebs_optimized = false
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
}

# Security Group for EFS
resource "aws_security_group" "sg_internal_efs" {
  name        = "engineering-tools-efs"
  description = "Allows Internal EFS use from Engineering VPC"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = ["${module.engineering-tools-ecs-cluster.security_group_id}"]
  }

  ingress {
    from_port = 2049
    to_port = 2049
    protocol = "tcp"
    security_groups = ["${aws_security_group.jenkins-slave.id}"]
  }

  egress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = ["${data.terraform_remote_state.engineering.cidr}"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "Engineering Tools EFS"
  }
}

module "rds-cluster" {
  source               = "../../modules/rds-cluster"
  master_username      = "lfengineering"
  name                 = "engineering-tools"
  master_password      = "buanCAWwwAGxUyoU2Fai"
  zone_id              = "${data.terraform_remote_state.engineering.zone_id}"
  availability_zones   = "${data.terraform_remote_state.engineering.availability_zones}"
  vpc_id               = "${data.terraform_remote_state.engineering.vpc_id}"
  subnet_ids           = "${data.terraform_remote_state.engineering.internal_subnets}"
  environment          = "Preprod"
  team                 = "Engineering"
  security_groups      = ["${module.engineering-tools-ecs-cluster.security_group_id}"]
  engine               = "mariadb"
  engine_version       = "10.1.19"
  parameter_group_name = "engineering"
  dns_name             = "rds"
  instance_type        = "db.t2.medium"
}
