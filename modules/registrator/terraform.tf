variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "region" {
  description = "The AWS Region in which we are providing the service in"
  default = "us-west-2"
}

variable "project" {
  description = "project name"
}

variable "log_group_name" {
  description = "Consul Server Endpoint the Agent needs to connect to"
}

data "template_file" "registrator_ecs_task" {
  template = "${file("${path.module}/registrator-ecs-task.json")}"

  vars {
    # For Logging
    AWS_REGION    = "${var.region}"
    PROJECT       = "${var.project}"
    LOG_GROUP_NAME = "${var.log_group_name}"
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
  desired_count                      = "20"
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
