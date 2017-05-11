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

variable "s3_bucket" {}

variable "redis_host" {}

variable "region" {}

variable "vpc_id" {}

data "template_file" "pypi_ecs_task" {
  template = "${file("${path.module}/pypi-ecs-task.json")}"

  vars {
    # Those are coming from the pypi-user on AWS, needs Read/Write on above S3 Bucket
    AWS_ACCESS_KEY_ID     = "AKIAIB6UWB7QG5QQYPWQ"
    AWS_SECRET_ACCESS_KEY = "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"
    S3_BUCKET             = "${var.s3_bucket}"
    REDIS_HOST            = "${var.redis_host}"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "pypi" {
  provider = "aws.local"
  family = "pypicloud"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.pypi_ecs_task.rendered}"
}

resource "aws_ecs_service" "pypi" {
  provider                           = "aws.local"
  name                               = "pypicloud-repository"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.pypi.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "pypi" {
  provider = "aws.local"
  name = "pypicloud-cluster"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 8080
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:8080"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "pypicloud-cluster"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "pypi" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.pypi.id}"
}

output "pypi_elb_cname" {
  value = "${aws_elb.pypi.dns_name}"
}

output "pypi_elb_name" {
  value = "${aws_elb.pypi.name}"
}

output "pypi_elb_zoneid" {
  value = "${aws_elb.pypi.zone_id}"
}

