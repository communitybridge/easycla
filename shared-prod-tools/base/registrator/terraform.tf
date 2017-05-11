variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

data "template_file" "registrator_ecs_task" {
  template = "${file("${path.module}/registrator-ecs-task.json")}"
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
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}
