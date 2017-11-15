variable "route53_zone" {}

variable "subdomain" {}

variable "project_name" {}

variable "instance_name" {}

variable "ecs_task_definition" {}

variable "container_port" {}

variable "container_name" {}

variable "public" {}

data "terraform_remote_state" "sandbox-env" {
  backend = "consul"
  config {
    address = "consul.service.development.consul:8500"
    path    = "terraform/sandboxes"
  }
}

resource "aws_ecs_service" "sandbox" {
  name                               = "${var.project_name}-${var.instance_name}"
  cluster                            = "${data.terraform_remote_state.sandbox-env.ecs_cluster_name}"
  task_definition                    = "${var.ecs_task_definition}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
  iam_role                           = "${data.terraform_remote_state.sandbox-env.ecs_role}"


  load_balancer {
    target_group_arn   = "${aws_alb_target_group.sandbox-main-port.arn}"
    container_name     = "${var.container_name}"
    container_port     = "${var.container_port}"
  }
}

resource "aws_alb_target_group" "sandbox-main-port" {
  name                 = "${var.project_name}-${var.instance_name}"
  port                 = "${var.container_port}"
  protocol             = "HTTP"
  vpc_id               = "${data.terraform_remote_state.sandbox-env.vpc_id}"
  deregistration_delay = 30

  health_check {
    path = "/"
    protocol = "HTTP"
    interval = 15
    unhealthy_threshold = 10
    healthy_threshold = 5
  }
}

resource "aws_alb" "sandbox-main" {
  name               = "${var.project_name}-${var.instance_name}"
  subnets            = ["${split(",", var.public == true ? join(",", data.terraform_remote_state.sandbox-env.external_subnets) : join(",", data.terraform_remote_state.sandbox-env.internal_subnets))}"]
  security_groups    = ["${var.public == true ? data.terraform_remote_state.sandbox-env.external_elb_sg : data.terraform_remote_state.sandbox-env.internal_elb_sg}"]
  internal           = "${var.public}"
}

resource "aws_alb_listener" "sandbox" {
  load_balancer_arn  = "${aws_alb.sandbox-main.id}"
  port               = "443"
  protocol           = "HTTPS"
  certificate_arn    = "${data.terraform_remote_state.sandbox-env.sandbox_cert_arn}"

  default_action {
    target_group_arn = "${aws_alb_target_group.sandbox-main-port.id}"
    type             = "forward"
  }
}

resource "aws_route53_record" "sandbox" {
  zone_id = "${var.route53_zone}"
  name    = "${var.subdomain}"
  type    = "A"

  alias {
    name                   = "${aws_alb.sandbox-main.dns_name}"
    zone_id                = "${aws_alb.sandbox-main.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "sandbox_www" {
  zone_id = "${var.route53_zone}"
  name    = "www.${var.subdomain}"
  type    = "A"

  alias {
    name                   = "${aws_alb.sandbox-main.dns_name}"
    zone_id                = "${aws_alb.sandbox-main.zone_id}"
    evaluate_target_health = true
  }
}
