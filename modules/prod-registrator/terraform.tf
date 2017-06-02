variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "dns_servers" {
  type = "list"
}

variable "region" {
  description = "The AWS Region in which we are providing the service in"
  default = "us-west-2"
}

variable "build_hash" {
  description = "The Build Number we are to deploy."
}

data "template_file" "registrator_ecs_task" {
  template = "${file("${path.module}/registrator-ecs-task.json")}"

  vars {
    # DNS Servers for Container Resolution
    DNS_SERVER_1  = "${var.dns_servers[0]}"
    DNS_SERVER_2  = "${var.dns_servers[1]}"
    DNS_SERVER_3  = "${var.dns_servers[2]}"

    # For Logging
    AWS_REGION    = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "registrator" {
  provider              = "aws.local"
  family                = "registrator"
  container_definitions = "${data.template_file.registrator_ecs_task.rendered}"
  network_mode          = "host"

  volume {
    name = "docker-sock"
    host_path = "/var/run/docker.sock"
  }
}

resource "aws_ecs_service" "registrator" {
  provider                           = "aws.local"
  name                               = "registrator"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.registrator.arn}"
  desired_count                      = "6"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  placement_strategy {
    type = "spread"
    field = "instanceId"
  }

  placement_constraints {
    type = "distinctInstance"
  }
}
