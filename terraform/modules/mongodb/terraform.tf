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

variable "route53_zone_id" {}

variable "data_path" {}

variable "vpc_id" {}

data "template_file" "mongodb_ecs_task" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "mongodb" {
  provider              = "aws.local"
  family                = "mongodb"
  container_definitions = "${data.template_file.mongodb_ecs_task.rendered}"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  volume {
    name      = "mongodb-storage"
    host_path = "${var.data_path}"
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

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "mongodb" {
  provider = "aws.local"
  name = "mongodb"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 27017
    instance_protocol = "tcp"
    lb_port = 27017
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:27017"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "mongodb"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "mongodb" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.mongodb.id}"
}

resource "aws_route53_record" "mongodb" {
  provider= "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "mongodb"
  type    = "A"

  alias {
    name                   = "${aws_elb.mongodb.dns_name}"
    zone_id                = "${aws_elb.mongodb.zone_id}"
    evaluate_target_health = true
  }
}
