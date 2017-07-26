variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "ecs_security_group" {}

variable "region" {}

variable "vpc_id" {}

variable "dns_servers" {
  type = "list"
}

variable "internal_subnets" {
  type = "list"
}

variable "newrelic_key" {}

variable "keypair" {}

variable "iam_role" {}

variable "region_identifier" {}

data "aws_ami" "amazon-linux-ecs-optimized" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["amzn-ami-*-amazon-ecs-optimized"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name = "owner-alias"
    values = ["amazon"]
  }
}

data "template_file" "ecs_instance_cloudinit_tools" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name  = "infrastructure"
    region            = "${var.region}"
    region_identifier = "${var.region_identifier}"
    newrelic_key      = "${var.newrelic_key}"
  }
}

resource "aws_instance" "ecs_mongodb" {
  ami = "${data.aws_ami.amazon-linux-ecs-optimized.id}"
  instance_type = "t2.small"
  subnet_id = "${var.internal_subnets[0]}"
  vpc_security_group_ids = ["${var.ecs_security_group}"]
  iam_instance_profile = "${var.iam_role}"
  key_name = "${var.keypair}"
  user_data = "${data.template_file.ecs_instance_cloudinit_tools.rendered}"

  tags {
    Name = "${var.ecs_cluster_name}-mongo"
  }
}

resource "aws_ebs_volume" "mongodb_storage" {
  availability_zone = "${var.region}a"
  size = 100
}

# Attach DBMS Volume
resource "aws_volume_attachment" "ebs_att" {
  device_name = "/dev/xvdh"
  volume_id = "${aws_ebs_volume.mongodb_storage.id}"
  instance_id = "${aws_instance.ecs_mongodb.id}"
}

data "template_file" "mongodb_ecs_task" {
  template = "${file("${path.module}/mongodb-ecs-task.json")}"

  vars {
    # DNS Servers for Container Resolution
    DNS_SERVER_1          = "${var.dns_servers[0]}"
    DNS_SERVER_2          = "${var.dns_servers[1]}"
    DNS_SERVER_3          = "${var.dns_servers[2]}"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "mongodb" {
  provider = "aws.local"
  family = "mongodb"
  network_mode = "host"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.mongodb_ecs_task.rendered}"

  volume {
    name = "storage"
    host_path = "/storage"
  }
}

resource "aws_ecs_service" "mongodb" {
  provider                           = "aws.local"
  name                               = "mongodb"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.mongodb.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  placement_constraints {
    expression = "attribute:ecs.instance-type == t2.small"
    type = "memberOf"
  }
}