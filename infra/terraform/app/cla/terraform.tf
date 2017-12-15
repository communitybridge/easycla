variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "internal_elb_sg" {
  description = "Security Group for the internal ELB"
}

variable "region" {
  description = "The AWS Region in which we are providing the service in"
  default = "us-west-2"
}

variable "vpc_id" {
  description = "The VPC ID we are deploying to"
}

variable "dns_servers" {
  type = "list"
}

variable "build_hash" {
  description = "The Build Number we are to deploy."
}

variable "ecs_role" {
  description = "The ecsService Role for cla"
}

variable "route53_zone_id" {
  description = "The Route53 Zone ID we need to add an entry to for the ALB/ELB."
}

data "template_file" "cla_ecs_task" {
  template = "${file("${path.module}/cla-ecs-task.json")}"

  vars {
    # Build Information
    build_hash              = "${var.build_hash}"

    # DNS Servers for Container Resolution
    DNS_SERVER_1              = "${var.dns_servers[0]}"
    DNS_SERVER_2              = "${var.dns_servers[1]}"
    DNS_SERVER_3              = "${var.dns_servers[2]}"

    # Tags for Registrator
    TAG_REGION                = "${var.region}"
    TAG_VPC_ID                = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "cla" {
  provider                = "aws.local"
  family                  = "cla"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions   = "${data.template_file.cla_ecs_task.rendered}"
}

resource "aws_ecs_service" "cla" {
  provider                           = "aws.local"
  name                               = "cla"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.cla.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
  iam_role                           = "${var.ecs_role}"

  load_balancer {
    target_group_arn   = "${aws_alb_target_group.cla.arn}"
    container_name     = "cla"
    container_port     = 5000
  }
}

resource "aws_alb_target_group" "cla" {
  provider             = "aws.local"
  name                 = "cla-5000"
  port                 = 5000
  protocol             = "HTTP"
  vpc_id               = "${var.vpc_id}"
  deregistration_delay = 30

  health_check {
    path = "/v1/health"
    protocol = "HTTP"
    interval = 15
    matcher = "200,202"
  }
}

resource "aws_alb" "cla" {
  provider           = "aws.local"
  name               = "cla"
  subnets            = ["${var.internal_subnets}"]
  security_groups    = ["${var.internal_elb_sg}"]
  internal           = true
}
