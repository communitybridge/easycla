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

variable "vpc_id" {}

variable "region" {}

variable "route53_zone_id" {}

variable "ecs_role" {}

data "template_file" "vault-ui_ecs_task" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "vault-ui" {
  provider              = "aws.local"
  family                = "vault-ui"
  container_definitions = "${data.template_file.vault-ui_ecs_task.rendered}"
}

resource "aws_ecs_service" "vault-ui" {
  provider                           = "aws.local"
  name                               = "vault-ui"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.vault-ui.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
  iam_role                           = "${var.ecs_role}"

  load_balancer {
    target_group_arn   = "${aws_alb_target_group.vault-ui.arn}"
    container_name     = "vault-ui"
    container_port     = 8080
  }
}

resource "aws_alb_target_group" "vault-ui" {
  provider             = "aws.local"
  name                 = "vault-ui-8080"
  port                 = 8080
  protocol             = "HTTP"
  vpc_id               = "${var.vpc_id}"
  deregistration_delay = 30

  health_check {
    path = "/"
    protocol = "HTTP"
    interval = 15
  }

  depends_on = [
    "aws_alb.vault-ui",
  ]
}

resource "aws_alb" "vault-ui" {
  provider           = "aws.local"
  name               = "vault-ui"
  subnets            = ["${var.internal_subnets}"]
  security_groups    = ["${var.internal_elb_sg}"]
  internal           = true
}

resource "aws_alb_listener" "vault-ui_80" {
  provider           = "aws.local"
  load_balancer_arn  = "${aws_alb.vault-ui.id}"
  port               = "80"
  protocol           = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.vault-ui.id}"
    type             = "forward"
  }
}

resource "aws_alb_listener" "vault-ui_443" {
  provider           = "aws.local"
  load_balancer_arn  = "${aws_alb.vault-ui.id}"
  port               = "443"
  protocol           = "HTTPS"
  certificate_arn    = "arn:aws:acm:us-west-2:643009352547:certificate/4938ed7c-e270-4597-84b2-6374db6149f4"

  default_action {
    target_group_arn = "${aws_alb_target_group.vault-ui.id}"
    type             = "forward"
  }
}

resource "aws_route53_record" "public" {
  provider= "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "vault-ui"
  type    = "A"

  alias {
    name                   = "${aws_alb.vault-ui.dns_name}"
    zone_id                = "${aws_alb.vault-ui.zone_id}"
    evaluate_target_health = true
  }
}
