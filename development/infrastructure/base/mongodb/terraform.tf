variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "region" {}

variable "vpc_id" {}

variable "dns_servers" {
  type = "list"
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
    host_path = "/mnt/storage/mongodb"
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
}