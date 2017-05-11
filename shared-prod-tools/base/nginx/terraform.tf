variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "ecs_asg_name" {
  description = "The name of the ECS AutoScaling Group"
}

variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "internal_elb_sg" {
  description = "Security Group for the internal ELB"
}

variable "region" {}

variable "vpc_id" {}

data "template_file" "nginx_ecs_task" {
  template = "${file("${path.module}/nginx-ecs-task.json")}"

  vars {
    TAG_REGION = "${var.region}"
    TAG_VPC_ID = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "nginx" {
  provider     = "aws.local"
  family       = "nginx-dns-lb"
  network_mode = "host"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.nginx_ecs_task.rendered}"
}

resource "aws_ecs_service" "nginx" {
  provider                           = "aws.local"
  name                               = "nginx"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.nginx.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}
