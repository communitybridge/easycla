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

variable "data_path" {}

variable "route53_zone_id" {}

data "template_file" "consul" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
//    DATA_PATH       = "${var.data_path}"
  }
}

resource "aws_ecs_task_definition" "consul" {
  provider = "aws.local"
  family   = "consul"
  network_mode          = "host"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.consul.rendered}"
}

resource "aws_ecs_service" "consul" {
  provider                           = "aws.local"
  name                               = "consul"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.consul.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "consul" {
  provider           = "aws.local"
  name               = "consul"
  security_groups    = ["${var.internal_elb_sg}"]
  subnets            = ["${var.internal_subnets}"]
  internal           = true

  listener {
    instance_port      = 8500
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:643009352547:certificate/4938ed7c-e270-4597-84b2-6374db6149f4"
  }

  listener {
    instance_port      = 8500
    instance_protocol  = "http"
    lb_port            = 80
    lb_protocol        = "http"
  }

  listener {
    instance_port      = 8301
    instance_protocol  = "tcp"
    lb_port            = 8301
    lb_protocol        = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:8300"
    interval            = 30
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Name = "consul"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "consul" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.consul.id}"
}

resource "aws_route53_record" "consul" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "consul"
  type    = "A"

  alias {
    name                   = "${aws_elb.consul.dns_name}"
    zone_id                = "${aws_elb.consul.zone_id}"
    evaluate_target_health = true
  }
}